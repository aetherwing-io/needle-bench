use crate::schema::{Schema, StructDef, Field, FieldType};

/// Generate Rust source code from a schema definition.
/// Produces struct definitions with getter methods.
pub fn generate(schema: &Schema) -> String {
    let mut output = String::new();
    output.push_str("// Auto-generated code — do not edit\n\n");

    for struct_def in &schema.structs {
        output.push_str(&generate_struct(struct_def));
        output.push('\n');
    }

    output
}

fn generate_struct(def: &StructDef) -> String {
    let mut out = String::new();

    // Struct definition
    out.push_str(&format!("pub struct {} {{\n", def.name));
    for field in &def.fields {
        let rust_type = type_to_rust(&field.field_type, false);
        out.push_str(&format!("    {}: {},\n", field.name, rust_type));
    }
    out.push_str("}\n\n");

    // Impl block with getters
    out.push_str(&format!("impl {} {{\n", def.name));
    for field in &def.fields {
        out.push_str(&generate_getter(field));
    }
    out.push_str("}\n");

    out
}

fn type_to_rust(ft: &FieldType, owned: bool) -> String {
    match ft {
        FieldType::String => {
            if owned { "String".to_string() } else { "String".to_string() }
        }
        FieldType::Int => "i64".to_string(),
        FieldType::Float => "f64".to_string(),
        FieldType::Bool => "bool".to_string(),
        FieldType::Ref(name) => name.clone(),
        FieldType::Optional(inner) => format!("Option<{}>", type_to_rust(inner, true)),
        FieldType::List(inner) => format!("Vec<{}>", type_to_rust(inner, true)),
    }
}

/// Generate a getter method for a field.
/// For simple types (int, float, bool), return by value.
/// For String and complex types, return a reference.
fn generate_getter(field: &Field) -> String {
    let name = &field.name;

    match &field.field_type {
        FieldType::Int | FieldType::Float | FieldType::Bool => {
            let rust_type = type_to_rust(&field.field_type, false);
            format!(
                "    pub fn {}(&self) -> {} {{\n        self.{}\n    }}\n\n",
                name, rust_type, name
            )
        }
        FieldType::String => {
            format!(
                "    pub fn {}(&self) -> &str {{\n        &self.{}\n    }}\n\n",
                name, name
            )
        }
        FieldType::Ref(type_name) => {
            // Ref types: return the nested struct via accessor
            format!(
                "    pub fn {}(&self) -> {} {{\n        &self.{}\n    }}\n\n",
                name, type_name, name
            )
        }
        FieldType::Optional(inner) => {
            let inner_ret = match inner.as_ref() {
                FieldType::String => "Option<&str>".to_string(),
                _ => format!("Option<&{}>", type_to_rust(inner, false)),
            };
            let body = match inner.as_ref() {
                FieldType::String => format!("self.{}.as_deref()", name),
                _ => format!("self.{}.as_ref()", name),
            };
            format!(
                "    pub fn {}(&self) -> {} {{\n        {}\n    }}\n\n",
                name, inner_ret, body
            )
        }
        FieldType::List(inner) => {
            let inner_type = type_to_rust(inner, false);
            format!(
                "    pub fn {}(&self) -> &[{}] {{\n        &self.{}\n    }}\n\n",
                name, inner_type, name
            )
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::schema;

    #[test]
    fn test_generate_simple_schema() {
        let schema = schema::load_builtin("simple").unwrap();
        let code = generate(&schema);
        assert!(code.contains("pub struct User"));
        assert!(code.contains("fn name(&self)"));
        assert!(code.contains("fn age(&self)"));
    }

    #[test]
    fn test_generate_complex_schema() {
        let schema = schema::load_builtin("complex").unwrap();
        let code = generate(&schema);
        assert!(code.contains("pub struct Employee"));
        assert!(code.contains("pub struct Company"));
        assert!(code.contains("pub struct Address"));
    }
}
