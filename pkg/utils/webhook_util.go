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
	"archive/tar"
	"bufio"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"gopkg.in/yaml.v2"
	"hash"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/klog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

/* constants */
const (
	TRIGGERS = "triggers"
	CHKSUM   = "sha256"
)

// ReadFile Reads a file and returns the bytes.
func ReadFile(fileName string) ([]byte, error) {
	ret := make([]byte, 0)
	file, err := os.Open(fileName)
	if err != nil {
		return ret, err
	}
	defer file.Close()
	input := bufio.NewScanner(file)
	for input.Scan() {
		for _, b := range input.Bytes() {
			ret = append(ret, b)
		}
	}
	return ret, nil
}

// ReadJSON reads a JSON file and returns it as an unstructured.Unstructured object.
func ReadJSON(fileName string) (*unstructured.Unstructured, error) {
	bytes, err := ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var unstructuredObj = &unstructured.Unstructured{}
	err = unstructuredObj.UnmarshalJSON(bytes)
	if err != nil {
		return nil, err
	}
	return unstructuredObj, nil
}

func downloadFileTo(url, path string) error {
	client := http.Client{}
	response, err := client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the file
	_, err = io.Copy(out, response.Body)
	return err
}

func getHTTPURLReaderCloser(url string) (io.ReadCloser, error) {

	client := http.Client{}
	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusOK {
		return response.Body, nil
	}
	return nil, fmt.Errorf("unable to read from url %s, http status: %s", url, response.Status)
}

// ReadHTTPURL Read remote file from URL and return bytes
func ReadHTTPURL(url string) ([]byte, error) {
	readCloser, err := getHTTPURLReaderCloser(url)
	if err != nil {
		return nil, err
	}
	defer readCloser.Close()
	bytes, err := ioutil.ReadAll(readCloser)
	return bytes, err
}

// YAMLToMap Converts a YAML byte array to a map.
func YAMLToMap(bytes []byte) (map[string]interface{}, error) {
	var myMap map[string]interface{}
	err := yaml.Unmarshal(bytes, &myMap)
	if err != nil {
		return nil, err
	}
	return myMap, nil
}

/*
DownloadTrigger Download the trigger.tar.gz and unpack into the directory
  triggerURL: URL that serves the trigger gzipped tar
  triggerChkSum: the sha256 checksum of the trigger archive
  dir: directory to unpack the trigger.tar.gz
*/
func DownloadTrigger(triggerURL, triggerChkSum, dir string, verifyChkSum bool) error {
	if klog.V(5) {
		klog.Infof("Entering downloadTrigger triggerURL: %s, directory to store trigger: %s", triggerURL, dir)
		defer klog.Infof("Leaving downloadTrigger triggerURL: %s, directory to store trigger: %s", triggerURL, dir)
	}

	triggerArchiveName := filepath.Join(dir, "incubator.trigger.tar.gz")
	err := downloadFileTo(triggerURL, triggerArchiveName)
	if err != nil {
		return err
	}

	// Verify that the checksum matches the value found in kabanero-index.yaml
	if verifyChkSum {
		chkSum, err := sha256sum(triggerArchiveName)
		if err != nil {
			return fmt.Errorf("unable to calculate checksum of file %s: %s", triggerArchiveName, err)
		}

		if klog.V(5) {
			klog.Infof("Calculated sha256 checksum of file %s: %s", triggerArchiveName, chkSum)
		}

		if chkSum != triggerChkSum {
			klog.Fatalf("trigger collection checksum does not match the checksum from the Kabanero index: found: %s, expected: %s",
				chkSum, triggerChkSum)
		}
	}

	// Untar the triggers collection
	triggerReadCloser, err := os.Open(triggerArchiveName)
	if err != nil {
		return err
	}

	err = DecompressGzipTar(triggerReadCloser, dir)
	return err
}

// MergePathWithErrorCheck Merge a directory path with a relative path. Return error if the rectory not a prefix of the merged path after the merge
func MergePathWithErrorCheck(dir string, toMerge string) (string, error) {
	dest := filepath.Join(dir, toMerge)
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	dest, err = filepath.Abs(dest)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(dest, dir) {
		return dest, nil
	}
	return dest, fmt.Errorf("unable to merge directory %s with %s, The merged directory %s is not in a subdirectory", dir, toMerge, dest)
}

/* Calculate the SHA256 sum of a file */
func sha256sum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

type HashFunc func() hash.Hash

func getHash(hashType string) (HashFunc, error) {
	switch hashType {
	case "sha1":
		return sha1.New, nil
	case "sha256":
		return sha256.New, nil
	}

	return nil, fmt.Errorf("unrecognized hash type '%s'", hashType)
}

func hashPayload(sigType, secret string, payload []byte) (string, error) {
	h, err := getHash(sigType)
	if err != nil {
		return "", err
	}

	hm := hmac.New(h, []byte(secret))
	hm.Write(payload)
	return hex.EncodeToString(hm.Sum(nil)), nil
}

// ValidatePayload verifies that a payload hashed with some secret matches the expected signature
func ValidatePayload(sigType, sigHash, secret string, payload []byte) error {
	hashedPayload, err := hashPayload(sigType, secret, payload)
	if err != nil {
		return err
	}

	if !hmac.Equal([]byte(hashedPayload), []byte(sigHash)) {
		return fmt.Errorf("payload hash is '%s' but expected '%s'", hashedPayload, sigHash)
	}

	klog.Infof("Payload validated with signature %s", hashedPayload)
	return nil
}

// DecompressGzipTar Decompresses and extracts a tar.gz file.
func DecompressGzipTar(readCloser io.ReadCloser, dir string) error {
	defer readCloser.Close()

	gzReader, err := gzip.NewReader(readCloser)
	if err != nil {
		return err
	}
	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if header == nil {
			continue
		}
		dest, err := MergePathWithErrorCheck(dir, header.Name)
		if err != nil {
			return err
		}
		fileInfo := header.FileInfo()
		mode := fileInfo.Mode()
		if mode.IsRegular() {
			fileToCreate, err := os.OpenFile(dest, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("unable to create file %s, error: %s", dest, err)
			}

			_, err = io.Copy(fileToCreate, tarReader)
			closeErr := fileToCreate.Close()
			if err != nil {
				return fmt.Errorf("unable to read file %s, error: %s", dest, err)
			}
			if closeErr != nil {
				return fmt.Errorf("unable to close file %s, error: %s", dest, closeErr)
			}
		} else if mode.IsDir() {
			err = os.MkdirAll(dest, 0755)
			if err != nil {
				return fmt.Errorf("unable to make directory %s, error:  %s", dest, err)
			}
			klog.Infof("Created subdirectory %s\n", dest)
		} else {
			return fmt.Errorf("unsupported file type within tar archive: file within tar: %s, field type: %v", header.Name, mode)
		}

	}
	return nil
}
