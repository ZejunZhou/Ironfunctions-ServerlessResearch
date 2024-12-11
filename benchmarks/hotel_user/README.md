```zsh
fn build
fn apps create hotel
fn routes delete hotel /user
fn routes create hotel /user -max-concurrency 8 -m 256 --format http --timeout 60s  --idle-timeout 600s
fn routes update hotel /user -max-concurrency 8 -m 256 --format http --timeout 60s  --idle-timeout 600s


curl -X POST --data '{
"username":"Cornell_1",
"password":"1111111111"
}' http://localhost:8080/r/hotel/user

REMOTE="pc99.cloudlab.umass.edu"

curl -X POST --data '{
"username":"Cornell_6802",
"password":"6802680268026802680268026802680268026802"
}' http://$REMOTE:8080/r/hotel/user

curl -X POST --data '{
"username":"Cornell_6802",
"password":"123123"
}' http://$REMOTE:8080/r/hotel/user
```
