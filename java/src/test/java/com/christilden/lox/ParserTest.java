package com.christilden.lox;

import java.util.Arrays;
import java.util.List;

import junit.framework.Test;
import junit.framework.TestCase;
import junit.framework.TestSuite;

/**
 * Unit test for AstPrinter.
 */
public class ParserTest extends TestCase {
    /**
     * Create the test case
     *
     * @param testName name of the test case
     */
    public ParserTest( String testName ) {
        super( testName );
    }

    /**
     * @return the suite of tests being tested
     */
    public static Test suite() {
        return new TestSuite( ParserTest.class );
    }

    /**
     * Tests the parser.
     */
    public void testParserSimple() {
        int line = 1;
        List<Token> tokens = Arrays.asList(new Token(TokenType.NUMBER, null, 123.0, line), new Token(TokenType.MINUS, "-", null, line), new Token(TokenType.NUMBER, null, 45.67, line), new Token(TokenType.EOF, "", null, line));

        Parser parser = new Parser(tokens);
        Expr expression = parser.parse();

        assertEquals(new AstPrinter().print(expression), "(- 123.0 45.67)");
    }

    /**
     * Tests the scanner and parser.
     */
    public void testScannerAndParserSimple() {
        Scanner scanner = new Scanner("1 + 2 * 3");
        List<Token> tokens = scanner.scanTokens();
        Parser parser = new Parser(tokens);
        Expr expression = parser.parse();

        assertEquals(new AstPrinter().print(expression), "(+ 1.0 (* 2.0 3.0))");
    }
}
