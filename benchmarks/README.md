## REMOTE

cloudlab: `pc70.cloudlab.umass.edu`

```zsh
curl -X POST --data '{
"username":"user123",
"password":"pass123"
}' http://pc70.cloudlab.umass.edu:8080/r/hotel/user
```

## REMOTE WRK2

```zsh
wrk -t 4 -c 64 -d 150 -L -U -s ./mixed-workload_type_1.lua http://pc70.cloudlab.umass.edu:8080 -R 10

```
