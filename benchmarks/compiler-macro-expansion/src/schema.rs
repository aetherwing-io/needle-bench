/// Represents a field in a generated struct.
#[derive(Debug, Clone)]
pub struct Field {
    pub name: String,
    pub field_type: FieldType,
}

/// Types that can appear in a schema.
#[derive(Debug, Clone)]
pub enum FieldType {
    String,
    Int,
    Float,
    Bool,
    /// A reference to another struct by name.
    Ref(String),
    /// An optional field.
    Optional(Box<FieldType>),
    /// A list of items.
    List(Box<FieldType>),
}

/// A struct definition in the schema.
#[derive(Debug, Clone)]
pub struct StructDef {
    pub name: String,
    pub fields: Vec<Field>,
}

/// A schema is a collection of struct definitions.
#[derive(Debug, Clone)]
pub struct Schema {
    pub name: String,
    pub structs: Vec<StructDef>,
}

/// Load a built-in schema for testing.
pub fn load_builtin(name: &str) -> Option<Schema> {
    match name {
        "simple" => Some(Schema {
            name: "simple".to_string(),
            structs: vec![
                StructDef {
                    name: "User".to_string(),
                    fields: vec![
                        Field { name: "name".to_string(), field_type: FieldType::String },
                        Field { name: "age".to_string(), field_type: FieldType::Int },
                        Field { name: "active".to_string(), field_type: FieldType::Bool },
                    ],
                },
            ],
        }),
        "complex" => Some(Schema {
            name: "complex".to_string(),
            structs: vec![
                StructDef {
                    name: "Address".to_string(),
                    fields: vec![
                        Field { name: "street".to_string(), field_type: FieldType::String },
                        Field { name: "city".to_string(), field_type: FieldType::String },
                        Field { name: "zip".to_string(), field_type: FieldType::String },
                    ],
                },
                StructDef {
                    name: "Company".to_string(),
                    fields: vec![
                        Field { name: "name".to_string(), field_type: FieldType::String },
                        Field { name: "address".to_string(), field_type: FieldType::Ref("Address".to_string()) },
                    ],
                },
                StructDef {
                    name: "Employee".to_string(),
                    fields: vec![
                        Field { name: "name".to_string(), field_type: FieldType::String },
                        Field { name: "email".to_string(), field_type: FieldType::Optional(Box::new(FieldType::String)) },
                        Field { name: "company".to_string(), field_type: FieldType::Ref("Company".to_string()) },
                        Field { name: "tags".to_string(), field_type: FieldType::List(Box::new(FieldType::String)) },
                    ],
                },
            ],
        }),
        _ => None,
    }
}
