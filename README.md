Jo is a collection of small tools that work with JSON data.

`jo do` invokes a command for each object in the json input with shell variable for each field in the object.


Example of `jo do`

Grep to show last 10 errors in all container logs:
```
ns=kube-system
kubectl -n $ns get po -o json | jq '[ .items[] | .metadata.name as $p | (.spec.containers + .spec.initContainers)[] | {"p": $p, "c": .name} ]' | jo do "echo '----' \$p \$c; kubectl -n $ns logs \$p -c \$c | grep -i error | tail -n 10"
```


Write all container logs to files:
Note that the command passed to jo is quoted (without the quotes all logs will go to a single file)
```
kubectl -n $ns get po -o json | jq '[ .items[] | .metadata.name as $p | (.spec.containers + .spec.initContainers)[] | {"p": $p, "c": .name} ]' | jo do "kubectl -n $ns logs --ignore-errors=true \$p -c \$c >$out_dir/\$p+\$c.log"
```


Most of the time jo do flags and command flag are separated ok. When flags are ambiguous add --
```
jo do --in example.json -- echo -n \$x
```


## Installation

```
go install github.com/mmlt/jo
```
(see $GOBIN or $HOME/go/bin)