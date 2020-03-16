# events-operator
[![Build Status](https://travis-ci.com/kabanero-io/events-operator.svg?branch=master)](https://travis-ci.com/kabanero-io/events-operator)

## Table of Contents
#[Introduction](#Introduction)
#[Functional Specification](#Functional_Spec)

<a name="Introduction"></a>
## Introduction

The events operator allows users to define a Kubernetes centric event mediation flow. Through custom resource definitions, users can quickly construct mediation logic to receive, transform, and route JSON data structure. 

<a name="Functional_Spec"></a>
## Functional Specification

The main components of events infrastructure are:
- event mediator: defines what is to be run within one container.  It it constis of an optional https listener, and a list of mediations.
- event mediation: user defined logic used to transform or route events.
- event connection: defines the flow of data between meidations.

Due to the used of CRDs, the mediators, mediations, and connections may be changed dynamically.

### Event Mediators

An event mediator contains a list of mediations. Here is an example:

```yaml
apiVersion: events.kabanero.io/v1alpha1
kind: EventMediator
metadata:
  name: webhook
spec:
  createListener: true
  createRoute: true
  mediations:
    - mediation:
        name: webhook
        input: message
        sendTo: [ "dest"  ]
        body:
          - = : 'sendEvent(dest, message.body, message.header)'
```

When the attribute `createListener` is `true`, a https listener is created to receive JSON data as input. 
In addition, a `Service` with the same name as the mediator's name is created so that the listener is acccessible. 
An Openshift service serving self-signed TLS certificate is automatically created to secure the communications. 
No authentication/authorization is currently implemented. 

The URL to send a JSON message to the mediation within the mediator is `https://<mediatorname>/<mediation name>`. 
For example: `https://webhook/webhook`.  
The `<mediation name>` in the URL addresses the specific mediation within the mediator.

When both attributes `createListener` and `createRoute` are set to `true`, a new `Route` with the same name as the mediator is created to allow external access to the mediator. 
The external host name for the `Route` is installation specific. 
The URL to send a message to the mediation is `https:<external name>/<mediator name>/<meidation name>`. 
For example: `https://webhook-default.apps.mycompany.com/webhook/webhook`.

### Event Mediators

Each event mediation within a mediator defines one path for message processing. 
Its general form looks like :

```yaml
  mediations:
    - mediation:
        name: <mediation name>
        input: <identifier for input message>
        sendTo: [ "destination 1", "destination 2", ...  ]
        body:
           <body>
```


The attributes are:
- name: the name of the mediation. Note that the URL to the meidator must include the mediation name as the component of the path.
- input: the name of the input variable that contains the input message.
- Sendto: list of variable names for destinations to send output emssage.
- body: body that contains code based on Common Expression Language (CEL) to process the message.

The `body` of a mediation is an array of JSON objects, where each object may contain one or multipels of:
- An assigmment
- An `if` statement
- A `switch` statement
- A `default` statement (if nested in a swtich statement)
- A nested `body`

Here are an examples:

```yaml
apiVersion: events.kabanero.io/v1alpha1
kind: EventMediator
metadata:
  name: example
spec:
  createListener: true
  createRoute: true
  mediations:
    - mediation:
        name: mediation1
        input: switchboard
        sendTo: [ "dest1", "dest2", "dest3"  ]
        body:
          - =: 'attrValue = "" '
          - if: "has(message.body.attr)"
            =: "attrValue = message.body.attr"
          - switch:
              - if : ' attrvalue == "value1" '
                =: "sendEvent(dest1, message.body, message.header)"
              - if : 'attrValue == "value2" '
                sendEvent(dest2, message.body, message.header)
              - default:
                =: "sendEvent(dest3, message.body, message.header)"
```

More formally, 
- A `body` is an array of JSON objects that may contain the attribute names : `=`, `if`, `switch`, and `default`.
- The valid combinations of the attribute names in the same JSON object are:
  - `=`: an single assignment statement 
  - `if` and `=` : The assignment is executed when the condition of the `if` is true
  - `if` and `body`: The body is executed when the condition of the if is true
  - `switch` and `body`: The body must be array of JSON objects, where each element of the array is either an `if` statement, or a `default` statement.

Here are examples of an assignments. Note that not using a variable is allowed.

```yaml
=: 'attrValue = 1"
=: " sendEvent(dest, message.body, message.header)
```

Here is the first variation of an `if` statement:

```yaml
 - if : ' attrvalue == "value1" '
   =: "sendEvent(dest1, message.body, message.header)"
```

And second variation of an `if` statement with a `body`:

```yaml
- if : ' attrvalue == "value1" '
  body:
    - =: "attr = "value1""
    - =: "sendEvent(dest1, message.body, message.header)"
```

Here is an example of `swtich` statement:

```yaml
- switch:
  - if : ' attrvalue == "value1" '
    =: "sendEvent(dest1, message.body, message.header)"
  - if : 'attrValue == "value2" '
    sendEvent(dest2, message.body, message.header)
  - default:
    =: "sendEvent(dest3, message.body, message.header)"
```

#### Build-in functions

##### filter

The filter function returns a new map or array with some elements of the original map or array filtered out.

Input:
- message: a map or array data structure
- conditional: CEL expression to evaluate each element of the data structure. If it evaluates to true, the element is kept in the returned data structure. Otherwise, it is discarded. For a map, the variable `key` is bound to the key of the element being evaluated, and the `value` variable is bound to the value. For an array, only the `value` variable is available.

Output: 
- A copy of the original data structure with some elements filtered out based on the condition.

Examples:

This example keeps only those elements of the input `header` variable that is set by github:

```yaml
 - newHader : ' filter(header, " key.startsWith(\"X-Github\") || key.startsWith(\"github\")) '
 ```


 This example keeps only those elements of an integer array whose value is less than 10:
```yaml
   - newArray: ' filter(oldArray, " value < 10 " )
```

##### call

The call function is used to call a user defined function.

input:
- name: name of the function
- param: parameter for the function

output:
- return value from the function


Example:

The function `sum` implements a recursive function to calculate sum of all numbers from 1 to input:
```yaml
functions:
  - name: sum
    input: input
    output: output
    body:
      - switch:
          - if : 'input <= 0'
            output : input
          - default:
            - output: ' input + call("sum", input- 1)'
```


##### sendEvent

The sendEvent function sends an event to a destination.

Input:
  - destination: destination to send the event
  - message: a JSON compatible message
  - context : optional context for the event, such as http header

Output: empty string if OK, otherwise, error message

Example:
```yaml
  - result: " sendEvent("tekton-listener", message,  header)
```


##### jobID

The jobID function returns a new unique string each time it is called.



##### toDomainName

The toDomainName function converts a string into domain name format.

Input: a string
Output: the string converted to domain name format 

##### toLabel


The toLabel function converts a string in to Kubernetes label format.

Input: a string
Output: the string converted to label format 

##### split

Split a string into an array of string

Input: 
  - str: string to split
  - separator: the separator to split
Output: array of string containing original string separated by the separator.

Example:
```yaml
  - components: " split('a/b/c', '/') "
```

After split, the variable components contains `[ "a", "b", "c" ]`.