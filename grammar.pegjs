{
  function extractDelimited(first, rest) {
    var res = [first];

    if (rest) {
      for (var i = 0, len = rest.length, thing; i < len; ++i) {
        thing = rest[i];
        for (var k = 0, len2 = thing.length, obj; k < len2; ++k) {
          obj = thing[k];
          if ("string" === typeof obj.kind) {
            res.push(obj);
          }
        }
      }
    }

    return res;
  }
}

Script = prog:Statement* EOF {
  return {
    kind: "root",
    children: prog
  };
}

Statement = EntityDef / EntityGen

EntityDef = _ "def" _ entity:Identifier _ '{' _ body:FieldSet? _ '}' _ {
  return {
    kind: "entity",
    name: entity.value,
    children: body || []
  };
}

EntityGen "entity generation" = _ "generate" _ name:Identifier _ args:Arguments _ {
  return {
    kind: "generation",
    name: name.value,
    args: args || []
  };
}

FieldSet = first:FieldDecl rest:(_ ',' _ FieldDecl)* (_ ',')? {
  return extractDelimited(first, rest);
}

StaticDecl = name:Identifier _ fieldValue:Literal _ {
  return {
    name: name.value,
    kind: "field",
    value: fieldValue
  };
}

SymbolicDecl = name:Identifier _ fieldType:Builtin _ args:Arguments? _ {
  return {
    name: name.value,
    kind: "field",
    value: fieldType,
    args: args || []
  };
}

FieldDecl = StaticDecl / SymbolicDecl

Literal = value:(DateTime / Number / Bool / StringLiteral) { return value; }

SingleArgument = value:(Literal / Identifier) { return value; }

Arguments "arguments" = '(' _ first:SingleArgument rest:(_ ',' _ SingleArgument)* _ ')' {
  return extractDelimited(first, rest);
}

Identifier = [a-z_]i[a-z0-9_]i* {
  return {
    kind: "variable",
    value: text()
  };
}

Builtin = ("integer" / "decimal" / "string" / "date" / "dict") {
  return {
    kind: "builtin",
    value: text()
  };
}

Number = '-'? INT ('.' DIGIT+)? {
  return {
    kind: "Literal",
    value: parseFloat(text())
  };
}

Bool = ("true" / "false") {
  return {
    kind: "Literal",
    value: text() === "true"
  }
}

// Supports ISO-8601 date, date with timestamp, and optional zone offset
DateTime = _ iso8601Date:(DIGIT DIGIT DIGIT DIGIT '-' DIGIT DIGIT '-' DIGIT DIGIT LocalTimePart?) _ {
  return {
    kind: "Literal",
    value: new Date(iso8601Date.join(""))
  };
}

LocalTimePart = 'T'i DIGIT DIGIT ':' DIGIT DIGIT ':' DIGIT DIGIT ZonePart? { return text(); }

ZonePart = ('Z'i / [+\-] DIGIT DIGIT ':' DIGIT DIGIT) { return text(); }

StringLiteral "string"= '"' chars:CHAR* '"' {
  return {
    kind: "Literal",
    value: chars.join("")
  };
}

CHAR = UNESCAPED / ESCAPED

ESCAPE = "\\"

UNESCAPED = [^\0-\x1F\x22\x5C]

ESCAPED = ESCAPE sequence:(LITERAL_SEQ / INVISIBLE_SEQ / UNICODE_SEQ) { return sequence; }

UNICODE_SEQ = 'u' digits:(HEXDIG HEXDIG HEXDIG HEXDIG) {
  return String.fromCharCode(parseInt(digits.join(""), 16));
}

INVISIBLE_SEQ =
      'b' { return "\b"; }
      / 'f' { return "\f"; }
      / 'n' { return "\n"; }
      / 'r' { return "\r"; }
      / 't' { return "\t"; }

LITERAL_SEQ = '"' / '\\' / '/'

INT = '0' / NON_ZERO DIGIT*

NON_ZERO = [1-9]

DIGIT = [0-9]

HEXDIG = [0-9a-f]i

_ "whitespace" = [ \t\r\n]*

EOF = !.
