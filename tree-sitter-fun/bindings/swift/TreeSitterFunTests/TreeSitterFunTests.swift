import XCTest
import SwiftTreeSitter
import TreeSitterFun

final class TreeSitterFunTests: XCTestCase {
    func testCanLoadGrammar() throws {
        let parser = Parser()
        let language = Language(language: tree_sitter_fun())
        XCTAssertNoThrow(try parser.setLanguage(language),
                         "Error loading Fun grammar")
    }
}
