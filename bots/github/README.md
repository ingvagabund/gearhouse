# Github bots

Automate routine operators over github PRs

## Usage

Via CLI:

```bash
$ export GITHUB_TOKEN=XXX
$ ./_output/bin/prlabeler --organization ORGANIZATION --repository REPOSITORY
```

Via Kubernetes:

```
$ kubectl create secret generic github-secret --from-literal=api_token=XXX
$ oc apply -f kubernetes/cronjob.yaml
```
