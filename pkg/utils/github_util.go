/*
Copyright 2020 IBM Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"context"
	"fmt"
    "strings"
	"github.com/google/go-github/github"
	// "k8s.io/client-go/kubernetes"
    "sigs.k8s.io/controller-runtime/pkg/client"
	"k8s.io/klog"
	"net/http"
)

/* Get the repository's information from from github message body: name, owner, html_url, and ref */
func getRepositoryInfo(body map[string]interface{}, repositoryEvent string) (string, string, string, string, error) {
	ref := ""

	if repositoryEvent == "push" {
		// use SHA for ref
		afterObj, ok := body["after"]
		if !ok {
			return "", "", "", "", fmt.Errorf("unable to find after for push webhook message")
		}
		ref, ok = afterObj.(string)
		if !ok {
			return "", "", "", "", fmt.Errorf("after for push webhook message (%s) is not a string but %T", afterObj, afterObj)
		}
	} else if repositoryEvent == "pull_request" {
		// use pull_request.head.sha
		prObj, ok := body["pull_request"]
		if !ok {
			return "", "", "", "", fmt.Errorf("unable to find pull_request in webhook message")
		}
		prMap, ok := prObj.(map[string]interface{})
		if !ok {
			return "", "", "", "", fmt.Errorf("pull_request in webhook message is of type %T, not map[string]interface{}", prObj)
		}
		headObj, ok := prMap["head"]
		if !ok {
			return "", "", "", "", fmt.Errorf("pull_request in webhook message does not contain head")
		}
		head, ok := headObj.(map[string]interface{})
		if !ok {
			return "", "", "", "", fmt.Errorf("pull_request head not map[string]interface{}, but is %T", headObj)
		}

		shaObj, ok := head["sha"]
		if !ok {
			return "", "", "", "", fmt.Errorf("pull_request.head.sha not found")
		}
		ref, ok = shaObj.(string)
		if !ok {
			return "", "", "", "", fmt.Errorf("pull_request merge_commit_sha in webhook message not a string: %v", shaObj)
		}
	}

	repositoryObj, ok := body["repository"]
	if !ok {
		return "", "", "", "", fmt.Errorf("unable to find repository in webhook message")
	}
	repository, ok := repositoryObj.(map[string]interface{})
	if !ok {
		return "", "", "", "", fmt.Errorf("webhook message repository object not map[string]interface{}: %v", repositoryObj)
	}

	nameObj, ok := repository["name"]
	if !ok {
		return "", "", "", "", fmt.Errorf("webhook message repository name not found")
	}
	name, ok := nameObj.(string)
	if !ok {
		return "", "", "", "", fmt.Errorf("webhook message repository name not a string: %v", nameObj)
	}

	ownerMapObj, ok := repository["owner"]
	if !ok {
		return "", "", "", "", fmt.Errorf("webhook message repository owner not found")
	}
	ownerMap, ok := ownerMapObj.(map[string]interface{})
	if !ok {
		return "", "", "", "", fmt.Errorf("webhook message repository owner object not map[string]interface{}: %v", ownerMapObj)
	}
	ownerObj, ok := ownerMap["login"]
	if !ok {
		return "", "", "", "", fmt.Errorf("webhook message repository owner login not found")
	}
	owner, ok := ownerObj.(string)
	if !ok {
		return "", "", "", "", fmt.Errorf("webhook message repository owner login not string : %v", ownerObj)
	}

	htmlURLObj, ok := repository["html_url"]
	if !ok {
		return "", "", "", "", fmt.Errorf("webhook message repository html_url not found")
	}
	htmlURL, ok := htmlURLObj.(string)
	if !ok {
		return "", "", "", "", fmt.Errorf("webhook message html_url not string: %v", htmlURL)
	}

	return owner, name, htmlURL, ref, nil
}

/*
DownloadYAML Downloads a YAML file from a git repository.
  kubeClient: controller client to API server
  namespace: namespace to look for secret to github
  secretName name of the secret containing the token to access github
  header: HTTP header from webhook
  bodyMap: HTTP  message body from webhook
*/
func DownloadYAML(kubeClient client.Client, namespace string, secretName string, header map[string][]string, bodyMap map[string]interface{}, fileName string) (map[string]interface{}, bool, error) {

	hostHeader, isEnterprise := header[http.CanonicalHeaderKey("x-github-enterprise-host")]
	var host string
	if !isEnterprise {
		host = "github.com"
	} else {
		host = hostHeader[0]
	}

	repositoryEvent := header["X-Github-Event"][0]

	owner, name, htmlURL, ref, err := getRepositoryInfo(bodyMap, repositoryEvent)
	if err != nil {
		return nil, false, fmt.Errorf("unable to get repository owner, name, or html_url from webhook message: %v", err)
	}

	user, token, err := GetGitHubSecret(kubeClient, namespace, secretName,  htmlURL)
	if err != nil {
		return nil, false, fmt.Errorf("unable to get user/token secret for URL %s: %v", htmlURL, err)
	}

	githubURL := "https://" + host

	bytes, found, err := DownloadFileFromGithub(owner, name, fileName, ref, githubURL, user, token, isEnterprise)
	if err != nil {
		return nil, found, err
	}
	retMap, err := YAMLToMap(bytes)
	return retMap, found, err
}

// DownloadFileFromGithub Downloads a file and returns: bytes of the file, true if file exists, and any error
func DownloadFileFromGithub(owner, repository, fileName, ref, githubURL, user, token string, isEnterprise bool) ([]byte, bool, error) {

//	if klog.V(5) {
		klog.Infof("downloadFileFromGithub owner: %v, repo: %v, file: %v, ref: %v, githubURL: %v, user: %v, isEnterprise: %v", owner, repository, fileName, ref, githubURL, user, isEnterprise)
//	}

	ctx := context.Background()

	tp := github.BasicAuthTransport{
		Username: user,
		Password: token,
	}
	/*
		tokenService := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tokenClient := oauth2.NewClient(ctx, tokenService)
	*/

	var err error
	var client *github.Client
	if isEnterprise {
		githubURL = githubURL + "/api/v3"
		client, err = github.NewEnterpriseClient(githubURL, githubURL, tp.Client())
		if err != nil {
			return nil, false, err
		}
	} else {
		client = github.NewClient(tp.Client())
	}

	var options *github.RepositoryContentGetOptions
	if ref != "" {
		options = &github.RepositoryContentGetOptions{Ref: ref}
	}

	/*
	       rc, err := client.Repositories.DownloadContents(ctx, owner, repository, fileName, options)
	       if err != nil {
	   		fmt.Printf("Error type: %T, value: %v\n", err, err)
	           return nil, false, err
	       }
	       defer rc.Close()
	   	buf, err := ioutil.ReadAll(rc)
	*/
	fileContent, _, resp, err := client.Repositories.GetContents(ctx, owner, repository, fileName, options)
    klog.Infof("downloadFileFromGithub status code: %v", resp.Response.StatusCode)
	if resp.Response.StatusCode == 200 {
		if fileContent != nil {
			if fileContent.Content == nil {
				return nil, true, fmt.Errorf("content for %v/%v/%v is nil", owner, repository, fileName)
			}

			content, err := fileContent.GetContent()
			if err != nil {
				klog.Infof("download File Form Github error %v", err)
			} else {
				klog.Infof("download File from Github: buffer %v", content)
			}
			return []byte(content), true, err
		}
		/* some other errors */
		return nil, false, fmt.Errorf("unable to download %v/%v/%v: not a file", owner, repository, fileName)
	} else if resp.Response.StatusCode == 404 {
		/* does not exist */
		return nil, false, nil
	} else {
		/* some other errors */
		return nil, false, fmt.Errorf("unable to download %v/%v/%v, http error %v", owner, repository, fileName, resp.Response.Status)
	}

}

/* Return true if header is from Github event */
func IsHeaderGithub(header map[string][]string) bool {

    _, ok := header["X-Github-Event"]
    return ok
}

func IsHeaderGithubEnterprise(header map[string][]string) bool {

    _, isEnterprise := header[http.CanonicalHeaderKey("x-github-enterprise-host")]
    return isEnterprise
}

/* Parse Github URL into server, org, and repo 
  Input: url for the github server, e.g., https://github.com/org/repo
  OUtput:
     server:  The server, e.g., github.com
     org:  The org portion of the url
     repo: The name of the repo
*/
func ParseGithubURL(url string) (server, org, repo string, err error) {
    url = strings.Trim(url, " ")
    index := strings.Index(url, "://")
    if index < 0 {
        return "", "", "", fmt.Errorf("Unable to parse url: %v", url)
    }
    prefix := url[0:index]
    if prefix != "http" && prefix != "https" &&  prefix != "HTTP" && prefix != "HTTPS"{
        return "", "", "", fmt.Errorf("Unable to parse url: %v", url)
    }
    if len(url) <= index+3 {
        return "", "", "", fmt.Errorf("Unable to parse url: %v", url)
    }
    remainder := url[index+3:]
    components := strings.Split(remainder, "/")
    if len(components) != 3 {
        return "", "", "", fmt.Errorf("Unable to parse url: %v", url)
    }
    server = components[0]
    org = components[1]
    repo = components[2]

    return server, org, repo, nil
}

