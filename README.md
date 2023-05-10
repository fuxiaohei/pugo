# PuGo

a simple static site generator.

#### Why ?

I wanted to learn how to use golang, so I decided to write a simple static site generator. Even 'Hugo' is fully static site generator, but I just want a simple blog generator.

## Usage

create a folder named 'demo' to start site project.

```bash
$ mkdir demo
$ cd demo
$ pugo init
```

build site:

```bash
$ pugo build ./cmd/pugo
```

site locally on 'http://localhost:18080/': 

```bash
$ pugo server
```

