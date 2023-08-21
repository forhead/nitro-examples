use aws_sdk_dynamodb::types::AttributeValue;
use aws_sdk_dynamodb::{Client, Error};

/// Query the table for an item matching the input values.
/// Returns true if the item is found; otherwise false.
pub async fn query_item(client: &Client, name: &String, table: &String) -> Result<(), Error> {
    let name_av = AttributeValue::S(name.to_string());

    let item = client
        .get_item()
        .table_name(table)
        .key("name", name_av)
        .send()
        .await?;

    dbg!(item.item);

    Ok(())
}
