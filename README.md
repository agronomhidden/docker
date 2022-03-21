# docker

Program runs docker container with following bash-command command and write output to AWS CloudWatch

```shell
go build  -o golang-test-task main.go
```

```shell
./golang-test-task --docker-image python --bash-command $'pip install pip -U && pip install tqdm && python -c \"import time\ncounter = 0\nwhile True:\n\tprint(counter)\n\tcounter = counter + 1\n\ttime.sleep
(0.1)\"' --aws-region eu-central-1 --aws-access-key-id KEY --aws-secret-access-key SECRET --cloudwatch-group test-group --cloudwatch-stream test-stream7
```