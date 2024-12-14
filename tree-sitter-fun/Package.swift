// swift-tools-version:5.3
import PackageDescription

let package = Package(
    name: "TreeSitterFun",
    products: [
        .library(name: "TreeSitterFun", targets: ["TreeSitterFun"]),
    ],
    dependencies: [
        .package(url: "https://github.com/ChimeHQ/SwiftTreeSitter", from: "0.8.0"),
    ],
    targets: [
        .target(
            name: "TreeSitterFun",
            dependencies: [],
            path: ".",
            sources: [
                "src/parser.c",
                // NOTE: if your language has an external scanner, add it here.
            ],
            resources: [
                .copy("queries")
            ],
            publicHeadersPath: "bindings/swift",
            cSettings: [.headerSearchPath("src")]
        ),
        .testTarget(
            name: "TreeSitterFunTests",
            dependencies: [
                "SwiftTreeSitter",
                "TreeSitterFun",
            ],
            path: "bindings/swift/TreeSitterFunTests"
        )
    ],
    cLanguageStandard: .c11
)
