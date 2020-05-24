npm install apidoc
npm install apidoc-markdown2
rm -rf www
rm README-API.md
./node_modules/.bin/apidoc -e ./node_modules -o ./www
./node_modules/.bin/apidoc-markdown2 -p ./www -o README-API.md
