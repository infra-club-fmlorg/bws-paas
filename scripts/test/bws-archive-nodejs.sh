cd $(dirname $0)
if [ ! -e application.zip ]; then
  zip -r application vite-project -x \*/.git/\* \*/node_modules/\*
fi
