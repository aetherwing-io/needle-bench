# compiler-macro-expansion

## Project

A Rust code generation engine (`codegen-engine`) that takes schema definitions and produces struct definitions with typed getter methods. Similar to how a proc macro or protobuf compiler generates accessor code from a schema. The engine handles primitive types, optional fields, lists, and references between structs.

## Symptoms

Simple schemas with flat structs (single-level field access) work correctly. When generating code for schemas with nested struct references (e.g., Employee -> Company -> Address), the generated getter methods have incorrect return type signatures. The simple test suite passes, but the complex nested accessor test fails. The generated code for reference-typed fields produces getters that claim to return one type but actually return another, which would cause lifetime and ownership issues in real compiled Rust code.

## Bug description

The code generator handles most field types correctly but makes an error when generating getters for fields that reference other structs. The generated method signature doesn't match what the method body actually returns. This is subtle because single-level access works fine due to Rust's auto-deref, but chained accessor patterns (`obj.a().b().c()`) would fail to compile or exhibit use-after-free. The fix requires understanding the difference between owned and borrowed return types in generated code.

## Difficulty

Hard

## Expected turns

8-12
