## Login to cluster
Use this when you need a fresh token.

~~~bash
oc login https://api.example:6443 --token=$TOKEN
~~~

## Get project list

~~~bash
oc get projects
~~~

## Fix stuck namespace
Delete finalizers and reapply.

~~~bash
oc get ns stuck -o json | jq '.spec.finalizers = []' | oc replace -f -
~~~
