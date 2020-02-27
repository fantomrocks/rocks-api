#!/bin/bash
#
# This will bundle the GraphQL schema tree into a single source file so it can be complied with the code base
#
PACKAGE="gqlschema"
FILE="bundle.go"
BASEFOLDER="$(dirname "$0")/../internal/graphql/schema"

# make the bundle file content compacting the *.graphql files
{
  echo "package $PACKAGE"
  echo ""
  echo "// GraphQL Schema Bundle; auto-created ", `date "+%F %R"`
  echo "const schema = \`"
  find "$BASEFOLDER/schema" -name '*.graphql' -print0 | xargs -0 -I{} sh -c "cat {}; echo ''"
  echo "\`"
} >"$BASEFOLDER/$FILE"
