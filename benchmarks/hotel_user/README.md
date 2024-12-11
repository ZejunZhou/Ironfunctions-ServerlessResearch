```zsh
fn build
fn apps create hotel
fn routes delete hotel /user
fn routes create hotel /user -max-concurrency 8 -m 256 --format http --timeout 60s  --idle-timeout 600s
fn routes update hotel /user -max-concurrency 8 -m 256 --format http --timeout 60s  --idle-timeout 600s


curl -X POST --data '{
"username":"user123",
"password":"pass123"
}' http://localhost:8080/r/hotel/user

curl -X POST --data '{
"username":"Cornell_1",
"password":"1111111111"
}' http://pc99.cloudlab.umass.edu:8080/r/hotel/user
```
