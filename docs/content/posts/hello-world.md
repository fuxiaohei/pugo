```toml
title = "Hello World"
slug = "hello-world"
description = "this is a demo post"
tags = ["hello"]
date = "2022-05-10 18:05:07"
template = "post.html"
draft = false
comment = true
author = "admin"
```

If you see this article, you have installed it successfully. Thank you for using PuGo. I hope you can use it happily.

<!--more-->

### Guide

- create a new post

```bash
pugo create post new-hello-world.md
```

`new-hello-world.md` is at `content/posts/new-hello-world.md`.

- create new page

```bash
pugo create page new-hello-page.md
```

`new-hello-page.md` is at `content/pages/new-hello-page.md`.

- build site

```bash
pugo build
```

all files are builded at `build/`.

- serve locally

```bash
pugo serve
```

http server is started at `localhost:18080`.

### Publish

read online docs:

- [deploy to github pages](#)
- [deploy to netlify](#)
- [deploy to self host](#)
