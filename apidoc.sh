npm install apidoc@0.22.1
npm install apidoc-markdown2@0.3.7
rm -rf www
rm README-API.md
./node_modules/.bin/apidoc -e ./node_modules -i ./ -o ./www
./node_modules/.bin/apidoc-markdown2 -p ./www -o README-API.md
