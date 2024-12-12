```zsh
REMOTE=pc21.cloudlab.umass.edu
curl -X POST --data '{
    "require": "dis",
    "lat": 123.23,
    "lon": 32.6
}' http://$REMOTE:8080/r/hotel/recommendations
```
