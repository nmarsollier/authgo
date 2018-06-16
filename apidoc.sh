rm -rf www
apidoc -o www
apidoc-markdown2 -p www -o README-API.md
