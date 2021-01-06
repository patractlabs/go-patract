# Dumper for erc20 token

Dump erc20 status to a bolt database and we can query the data:

```bash
./dump -hash 5DTSQx8MnnMUHbYqWGn4wtRR4rgCTsmYdj6DKpLkbrVabFv7 -metadata ../../test/contracts/ink/erc20.json
```

Query data:

```bash
 curl http://127.0.0.1:8899/get/5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY
{"sequence":1,"address":"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY","amount":"10000000000000000000000000"}%
```
