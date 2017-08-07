;(function() {
  "use strict";
  function rootNode(statements) {
    return {
      kind: "root",
      children: searchNodes(statements)
    };
  }

  function importNode(path) {
    return {
      kind: "import",
      value: path
    };
  }

  function genNode(entity, args) {
    return {
      kind: "generation",
      value: entity,
      args: args || []
    }
  }

  function assignNode(identNode) {
    return {
      kind: "assignment",
      name: identNode.value,
    }
  }

  function entityNode(assignment, entity) {
    if (assignment) {
      entity.name = assignment.name
    }
    return entity;
  }

  function entityDefNode(extended, body) {
    var node = {
      kind: "entity",
      children: body || []
    };

    if (extended) {
      node.related = extended
    }

    return node;
  }

  function staticFieldNode(ident, fieldValue) {
    return {
      kind: "field",
      name: ident.value,
      value: fieldValue
    };
  }

  function dynamicFieldNode(ident, fieldType, args, bound) {
    return {
      kind:  "field",
      name:  ident.value,
      value: fieldType,
      args:  args || [],
      bound: bound || []
    };
  }

  function idNode(name) {
    return {
      kind: "identifier",
      value: name
    };
  }

  function builtinNode(value) {
    return {
      kind: "builtin",
      value: value
    };
  }

  function dateLiteralNode(date, localTime) {
    if (!!localTime && !!localTime) {
      date += localTime;
    }

    return {
      kind:  "literal-date",
      value: new Date(date)
    };
  }

  function floatLiteralNode(s) {
    return {
      kind:  "literal-float",
      value: parseFloat(s)
    };
  }

  function intLiteralNode(s) {
    return {
      kind:  "literal-int",
      value: parseInt(s, 10)
    };
  }
  function boolLiteralNode(value) {
    return {
      kind:  "literal-bool",
      value: "true" === value.toLowerCase()
    };
  }

  function strLiteralNode(value) {
    return {
      kind:  "literal-string",
      value: value
    };
  }

  function nullLiteralNode() {
    return {
      kind: "literal-null",
      value: null
    }
  }

  function searchNodes(v) {
    if (!v || (Array.isArray(v) && !v.length)) return [];
    if (v && "string" === typeof v.kind) return [v];

    for (var i = 0, res = [], cur, len = v.length; i < len; ++i) {
      cur = v[i];

      if (cur && "string" === typeof cur.kind) {
        res.push(cur);
      } else {
        if (Array.isArray(cur)) {
          res = res.concat(searchNodes(cur));
        }
      }
    }

    return res;
  }

  function delimitedNodeSlice(first, rest) {
    var res = [first];

    if (rest) {
      res = res.concat(searchNodes(rest));
    }

    return res;
  }

  module.exports = {
    rootNode,
    importNode,
    genNode,
    assignNode,
    entityNode,
    entityDefNode,
    staticFieldNode,
    dynamicFieldNode,
    idNode,
    builtinNode,
    dateLiteralNode,
    floatLiteralNode,
    intLiteralNode,
    boolLiteralNode,
    strLiteralNode,
    nullLiteralNode,
    searchNodes,
    delimitedNodeSlice
  };
})();
