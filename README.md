# craftinginterpreters
Building an interpreter and compiler!

# Dependencies
1. java 1.7+
2. maven 3.3.9+
3. python 2.7 (optional)

# Running (manually)

`cd java`

`mvn compile`

`mvn package`

Start a REPL to test the lox language (http://craftinginterpreters.com/the-lox-language.html)

`java -jar target/java-interp-1.0-SNAPSHOT.jar`

# Running (via Python's Fabric)

`cd java`

One time setup: `pip install -r requirements.txt`

Compile, package, and run the lox REPL

`fab run`
