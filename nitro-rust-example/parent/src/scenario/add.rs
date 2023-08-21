/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0.
 */

use aws_sdk_dynamodb::types::{AttributeValue, ReturnValue};
use aws_sdk_dynamodb::{Client, Error};

/// For add_item and query_item
pub struct Item {
    pub name: String,
    pub key_id: String,
    pub encrypted_private_key: String,
    pub address: String,
    pub encrypted_data_key: String,
}

#[derive(Debug, PartialEq)]
pub struct ItemOut {
    pub name: Option<AttributeValue>,
    pub key_id: Option<AttributeValue>,
    pub encrypted_private_key: Option<AttributeValue>,
    pub address: Option<AttributeValue>,
    pub encrypted_data_key: Option<AttributeValue>,
}

// Add an item to a table.
// snippet-start:[dynamodb.rust.add-item]
pub async fn add_item(client: &Client, item: Item, table: &String) -> Result<(), Error> {
    let name_av = AttributeValue::S(item.name);
    let key_id_av = AttributeValue::S(item.key_id);
    let encrypted_private_key_av = AttributeValue::S(item.encrypted_private_key);
    let address_av = AttributeValue::S(item.address);
    let encrypted_data_key_av = AttributeValue::S(item.encrypted_data_key);

    let request = client
        .put_item()
        .table_name(table)
        .item("name", name_av)
        .item("key_id", key_id_av)
        .item("encryptedPrivateKey", encrypted_private_key_av)
        .item("address", address_av)
        .item("encryptedDataKey", encrypted_data_key_av)
        .return_values(ReturnValue::None);

    // .return_values(ReturnValue::AllOld);
    // println!("Executing request [{request:?}] to add item...");

    let _resp = request.send().await?;

    // let attributes = resp.attributes().unwrap();
    // let key_id = attributes.get("KeyId").cloned();
    // let name = attributes.get("name").cloned();
    // println!("Added item {:?} {:?}", key_id, name);
    // Ok(ItemOut { key_id, name })

    Ok(())
}
