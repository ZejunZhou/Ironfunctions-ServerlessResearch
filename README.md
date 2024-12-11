## Get Start

1. `docker build -f Dockerfile.2690 -t 2690 .`

2. `docker run -d --rm -it --name functions -v ${PWD}/data:/app/data -v /var/run/docker.sock:/var/run/docker.sock -p 8080:8080 2690`

3. `cd test-function`

4. touch function.go

5. add code

Example

```
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Person struct {
	Name string
}

func main() {
	p := &Person{Name: "World"}
	json.NewDecoder(os.Stdin).Decode(p)
	fmt.Printf("Hello %v!", p.Name)
}
```

6. Install CLI tools

For Linux

```zsh
curl -LSs git.io/ironfn | sh
```

For macOS

```zsh
brew install iron-functions
```

7. Init function with runtime

```zsh
fn init --runtime go hello
```

8. Buld function - fn build is used to build the Docker image of the function

```zsh
fn build
```

9. Create app - Create application, organization and management multiple functions

```zsh
fn apps create myapp
```

10. Creates a route and maps the specified path to the newly created function

```zsh
fn routes create myapp /hello
```

11. If update function.go, then build and update route again

```zsh
fn build
fn routes update myapp /hello
```

12. Invoke function

```zsh
curl http://localhost:8080/r/myapp/hello
```

13. Frontend Access

```zsh
docker run -d --rm -it --link functions:api -p 4000:4000 -e "API_URL=http://api:8080" iron/functions-ui:0.0.2
```

Access http://localhost:4000/
