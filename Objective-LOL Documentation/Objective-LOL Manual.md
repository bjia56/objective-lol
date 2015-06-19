# Objective-LOL Language Specifications
---
---
<b>
Language Version: 0.1.0  
Compiler Version: 0.0.1
</b>
## (0) Contents
---

Use the search function on your viewer to jump to each labeled section.

<IMPORT FROM CONTENTS FILE>

## (1) Introduction
---

Objective-LOL is a variation of the esoteric programming language LOLCODE, created by Adam Lindsay, that is designed to support a functional and object-oriented programming model. As the syntax of LOLCODE relies heavily on the speech patterns of lolcat memes, Objective-LOL attempts to adopt a similar semantical model. Although being a language inspired by LOLCODE, Objective-LOL does not attempt to adhere to the language rules that drive LOLCODE; as such, a LOLCODE program may not compile without editing by the programmer.

Following the guidelines of other programming languages such as C++ and Java, a standalone, runnable Objective-LOL program requires that there is a global `MAIN` function, defined as the start of the program. If standalone runnability is not a requirement, the compiler may be instructed to produce a library containing the Objective-LOL code.

## (2) Language Fundamentals
---
The syntax of the Objective-LOL language and related keywords derive from the speech of lolcats. As such, some statements in Objective-LOL may not directly translate lexically to standard equivalents in other programming languages. This section attempts to explain the valid expressions of Objective-LOL, as well as how to combine them into meaningful statements.

### (2.a) Lexical Conventions
Identifiers in Objective-LOL exist as strings of any combination of letters, digits, and underscores. Similar to many other languages, Objective-LOL does not the use of a leading digit in identifiers.

The suggested naming convention for Objective-LOL is to adopt the lolcat convention of all uppercase text. This is not required, but highly suggested. Once compiled, however, all identifiers will be translated into uppercase text. Therefore, identifiers in Objective-LOL are not case-sensitive.

To follow the pattern of uppercase text, all operators and keywords are represented in uppercase. The following is a list of all keyword phrases and operators present, sorted in alphabetical order. These phrases are reserved by the language and should not be used as identifiers. Note that some keywords contain more than one string. Each listed phrase is accompanied by a short description of what it is used for.

### (2.b) Keywords
The following list is a collection of all reserved words and symbols pertinent to the Objective-LOL language, as well as each words' functionality and use.

#### (2.b.1) AN
Used in logic expressions. Returns `YEZ` if the two expressions on either side of this operator evaluate to `YEZ`, `NO` otherwise.

An example of an `IZ` statement utilizing `AN`:

    IZ YEZ AN NO?
        BTW this code will not run KK
    KTHX

When used in conjunction with the keyword `WIT`, `AN` denotes a continuation in the listing of arguments in function declarations:

    DIS TEH FUNCSHUN MAHFUNC WIT ARG1 TEH BOOL AN WIT ARG2 TEH INTEGR
        BTW some code here KK
    KTHX

Similarly, `AN WIT` is used to pass multiple arguments into a function call:

    MAX WIT 4 AN WIT 5

#### (2.b.2) AS
Used to explicitly cast a value to another type. An explicit cast is necessary when transforming a parent type to a child type. `AS` can also be used to cast primitives or from a child to a parent, although such an operation can be done implicitly, without the `AS` keyword.

Assuming that the class `CAT` is the parent class of `KITTEN`, the following is an example of using `AS` to cast:

    I HAS A KITTENVAR TEH KITTEN ITZ NEW KITTEN
    I HAS A CATVAR TEH CAT ITZ KITTENVAR AS CAT
    KITTENVAR ITZ CATVAR AS KITTEN

Objects can be cast to types not in their class hierarchies if the object to cast from has implemented a custom `AS` operator to cast to the target type.

#### (2.b.3) BIGGR THAN
Used in logic expressions. Requires the expressions on either side of this operator to be `NUMBR` values, or the left expression to have implemented a custom `BIGGR THAN` operator that takes the right expression as the singular argument. Returns `YEZ` if the expression on the left is numerically greater than the expression on the right.

An example of an `IZ` statement utilizing `BIGGR THAN`:

    IZ 5 BIGGR THAN 4?
        BTW this code will run KK
    KTHX

#### (2.b.4) BOOL
Used to explicitly declare a variable as a boolean value, or to declare that a function has a boolean return type.

The two possible values of any `BOOL` expression are:

    YEZ (true)
    NO (false)

These two values also act as constants to denote `BOOL` states.

If an `INTEGR` value is used as a `BOOL` value, a `NO` will be substituted as the `INTEGR` if the `INTEGR` value is zero; otherwise, a `YEZ` will be used. `DUBBLE` values cannot be used as a `BOOL` unless converted into an `INTEGR`.

Examples of a local `BOOL` variable declarations:

    I HAS A BOOLVAR TEH BOOL ITZ YEZ
    I HAS A BOOLVAR2 TEH BOOL ITZ 5

#### (2.b.5) BTW
Used to define the start of a comment block. All text, including text spanning multiple lines, held within a comment will be ignored when executing the program. Comments are closed with the `KK` keyword.

Example of comments:

    BTW the text here will not run KK
    I HAS A BOOLVAR BTW inline as well KK TEH BOOL ITZ YEZ
    BTW start of a multiline
        continuation of a multiline
        end of a multiline KK

#### (2.b.6) CLAS
Used to declare a class. Currently, classes can only be declared with the global scope.

An example of a class declaration:

    HAI ME TEH CLAS MAHCLAS
        BTW some code here KK
    KTHXBAI

A class can inherit from a parent class with the `KITTEH OF` expression:

    HAI ME TEH CLAS CAT
        BTW some code here KK
        BTW all EVRYONE and MAHSELF variables and functions can be accessed by children KK
    KTHXBAI

    HAI ME TEH CLAS KITTEN TEH KITTEH OF CAT


#### (2.b.8) DIS TEH
Used to declare a variable or a function with class scope. Any declarations with `DIS TEH` cannot be enclosed inside a function or outside of a class. Function declarations with `DIS TEH` must be closed by `KTHX`.

An example of a variable declaration inside a class:

    HAI ME TEH CLAS MAHCLAS
        DIS TEH VARIABLE INTVAR TEH INTEGR
    KTHXBAI

An example of a function declaration inside a class:

    HAI ME TEH CLAS MAHCLAS
        DIS TEH FUNCSHUN MAHFUNC TEH INTEGR
            BTW some code here KK
        KTHX
    KTHXBAI

#### (2.b.9) DIVIDEZ

#### (2.b.10) DUBBLE
Used to explicitly declare a variable as a double-precision value, or to declare that a function has a double-precision return type. How these values are stored in physical memory is determined by the virtual machine.

`DUBBLE` constants may be declared in code by either postfixing `D` to a number, or by using a decimal point `.` to denote a fractional part of the number.

Any `INTEGR` or `BOOL` value may be assigned to a `DUBBLE`, with automatic conversion. `INTEGR` values will retain the original values, while `BOOL` values will convert to `0D` for `NO` and `1D` for `YEZ`.

Examples of local `DUBBLE` variable declarations:

    I HAS A DUBBLEVAR TEH DUBBLE ITZ 5D
    I HAS A DUBBLEVAR2 TEH DUBBLE ITZ 2.5
    I HAS A DUBBLEVAR3 TEH DUBBLE ITZ 10
    I HAS A DUBBLEVAR4 TEH DUBBLE ITZ YEZ

#### (2.b.11) EVRYONE
Used to declare visibility of a variable or a function inside of a class. This keyword's visibility is the equivalent of `public` in other languages.

Declaring the visibility of variables or functions in a class is done by preceeding a section of declarations with the visibility term. An example is as follows:

    HAI ME TEH CLAS MAHCLAS
    EVRYONE
        DIS TEH VARIABLE MAHINT TEH INTEGR
        DIS TEH VARIABLE MAHSTR TEH STRIN
    KTHXBAI

#### (2.b.12) FUNCSHUN
Used to declare a function. Functions can be declared with global scope or class scope.

An example of a function declaration:

    HAI ME TEH FUNCSHUN MAHFUNC TEH STRIN WIT ARGS1
        BTW some code here KK
    KTHXBAI

#### (2.b.13) GIVEZ
Used to return a value from a function. Functions declared with a return type must return a value of that type or `NOTHIN`. Functions declared without a return type may use `GIVEZ UP` to exit early.

An example of returning a value from a function:

    HAI ME TEH FUNCSHUN GIVEZBOOL TEH BOOL
        GIVEZ YEZ
    KTHXBAI

An example of exiting a function without a return type:

    HAI ME TEH FUNCSHUN EXITEARLY
        WHILE YEZ
            IZ 5 BIGGR THAN 4?
                GIVEZ UP
            KTHX
        KTHX
    KTHXBAI

#### (2.b.14) GIVEZ UP

#### (2.b.15) HAI ME
Used to declare a variable, function, or class with global scope. Any declarations with `HAI ME` cannot be enclosed inside a function or a class, and must be closed with `KTHXBAI`.

An example of a global variable declaration:

    HAI ME TEH VARIABLE GLOBALINTVAR TEH INTEGR

An example of a global function declaration:

    HAI ME TEH FUNCSHUN GLOBALFUNC TEH INTEGR
        BTW some code here KK
    KTHXBAI

#### (2.b.16) I CAN HAS
Used to declare the libraries used by the current file. Equivalent to `#include` in C++ and `import` in Java. Exists to efficiently choose what Objective-LOL libraries are required and load those into memory. Lines loading libraries must be placed at the beginning of the file. The end of an import is optionally closed by a question mark `?`.

Examples of library loading:

    I CAN HAS STDIO?
    I CAN HAS MATH

If a `?` is used to close an import line, further library loads can be chained on the same line without writing `I CAN HAS`. Examples:

    I CAN HAS STDIO? MATH
    I CAN HAS STIDO? LIST? MATH? STRMANIP?

To load local files, the library name is replaced by a quoted string representing the relative file path to the current file. Example:

    I CAN HAS "otherfile.lol"?

#### (2.b.17) I HAS A
Used to declare a variable with local scope. Any declarations with `I HAS A` cannot be ouside of a function.

An example of a local variable declaration inside a function:

    HAI ME TEH FUNCSHUN MAHFUNC
        I HAS A INTVAR TEH INTEGR
    KTHXBAI

#### (2.b.18) IN
Used to access member functions and variables
