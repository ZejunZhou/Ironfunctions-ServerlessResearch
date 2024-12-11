## REMOTE

cloudlab: `pc70.cloudlab.umass.edu`

```zsh
curl -X POST --data '{
"username":"user123",
"password":"pass123"
}' http://pc70.cloudlab.umass.edu:8080/r/hotel/user
```

## Task Testing

### HOTEL_USER

```zsh
wrk -L -U -s ./hotel_user.lua http://pc99.cloudlab.umass.edu:8080 -R 50 > ./hotel_user/optimized/50
wrk -L -U -s ./hotel_user.lua http://pc99.cloudlab.umass.edu:8080 -R 100 > ./hotel_user/optimized/100
wrk -L -U -s ./hotel_user.lua http://pc99.cloudlab.umass.edu:8080 -R 200 > ./hotel_user/optimized/200
wrk -L -U -s ./hotel_user.lua http://pc99.cloudlab.umass.edu:8080 -R 300 > ./hotel_user/optimized/300
wrk -L -U -s ./hotel_user.lua http://pc99.cloudlab.umass.edu:8080 -R 400 > ./hotel_user/optimized/400
wrk -L -U -s ./hotel_user.lua http://pc99.cloudlab.umass.edu:8080 -R 500 > ./hotel_user/optimized/500
```

nightcore:

```zsh
curl http://pc99.cloudlab.umass.edu:8080/function/user?username=Cornel_2&password=2222222222

wrk -L -U -s ./hotel_user_nightcore.lua http://pc99.cloudlab.umass.edu:8080 -R 50 > ./hotel_user/nightcore/50
wrk -L -U -s ./hotel_user_nightcore.lua http://pc99.cloudlab.umass.edu:8080 -R 100 > ./hotel_user/nightcore/100
wrk -L -U -s ./hotel_user_nightcore.lua http://pc99.cloudlab.umass.edu:8080 -R 200 > ./hotel_user/nightcore/200
wrk -L -U -s ./hotel_user_nightcore.lua http://pc99.cloudlab.umass.edu:8080 -R 300 > ./hotel_user/nightcore/300
wrk -L -U -s ./hotel_user_nightcore.lua http://pc99.cloudlab.umass.edu:8080 -R 400 > ./hotel_user/nightcore/400
wrk -L -U -s ./hotel_user_nightcore.lua http://pc99.cloudlab.umass.edu:8080 -R 500 > ./hotel_user/nightcore/500
```

### Mixed Type

## REMOTE WRK2

```zsh
wrk -t 4 -c 64 -d 150 -L -U -s ./mixed-workload_type_1.lua http://pc70.cloudlab.umass.edu:8080 -R 10

```
