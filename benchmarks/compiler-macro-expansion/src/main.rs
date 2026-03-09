use std::env;
use std::process;

mod codegen;
mod schema;
mod runtime;

fn main() {
    let args: Vec<String> = env::args().collect();

    match args.get(1).map(|s| s.as_str()) {
        Some("generate") => {
            let schema_name = args.get(2).unwrap_or_else(|| {
                eprintln!("Usage: codegen-engine generate <schema-name>");
                process::exit(1);
            });
            match schema::load_builtin(schema_name) {
                Some(s) => {
                    let code = codegen::generate(&s);
                    println!("{}", code);
                }
                None => {
                    eprintln!("Unknown schema: {}", schema_name);
                    process::exit(1);
                }
            }
        }
        Some("test-simple") => {
            println!("Running simple accessor tests...");
            let result = runtime::test_simple();
            if result {
                println!("PASS: simple accessor tests");
            } else {
                println!("FAIL: simple accessor tests");
                process::exit(1);
            }
        }
        Some("test-complex") => {
            println!("Running complex nested accessor tests...");
            let result = runtime::test_complex();
            if result {
                println!("PASS: complex nested accessor tests");
            } else {
                println!("FAIL: complex nested accessor tests");
                process::exit(1);
            }
        }
        _ => {
            eprintln!("Usage: codegen-engine <generate|test-simple|test-complex>");
            process::exit(1);
        }
    }
}
