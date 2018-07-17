# etcd-vault-transformer

An [`etcdctl migrate` transformer](https://github.com/coreos/etcd/tree/master/etcdctl#migrate-options)
to help migrate a [vault](https://www.vaultproject.io/) instance using an
[etcd v2 API storage backend](https://www.vaultproject.io/docs/configuration/storage/etcd.html) to etcd v3 API.
i.e. change `etcd_api = "v2"` to `etcd_api = "v3"`

## Build

```
glide install
go build -o etcd-vault-transformer main.go
```

## Usage

```
etcdctl migrate --data-dir=/var/etcd --transformer=etcd-vault-transformer
```

## References

 * https://github.com/hashicorp/vault/blob/v0.10.1/physical/etcd/etcd2.go
 * https://github.com/hashicorp/vault/issues/2862
 * https://github.com/hongchaodeng/etcdupgrade/blob/master/main.go
