```zsh
fn build

fn apps create myapp
fn routes create hotel /rate -max-concurrency 8 -m 512 --format http --timeout 60s --idle-timeout 600s
fn routes update hotel /rate -max-concurrency 8 -m 512 --format http --timeout 60s --idle-timeout 600s
```

```zsh
curl http://localhost:8080/r/hotel/rate

fn routes update hotel /rate -max-concurrency 8 -m 512 --format http --timeout 60s --idle-timeout 600s
```
