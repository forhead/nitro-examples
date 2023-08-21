use aws_config::{meta::region::RegionProviderChain, SdkConfig};
use aws_sdk_dynamodb::{config::Region, meta::PKG_VERSION, Error};

#[derive(Debug)]
pub struct BaseOpt {
    /// The AWS Region.
    pub region: Option<String>,

    /// Whether to display additional information.
    pub verbose: bool,
}

pub fn make_region_provider(region: Option<String>) -> RegionProviderChain {
    RegionProviderChain::first_try(region.map(Region::new))
        .or_default_provider()
        .or_else(Region::new("ap-east-1"))
}

pub async fn make_config(opt: BaseOpt) -> Result<SdkConfig, Error> {
    let region_provider = make_region_provider(opt.region);

    if opt.verbose {
        println!("DynamoDB client version: {}", PKG_VERSION);
        println!(
            "Region:                  {}",
            region_provider.region().await.unwrap().as_ref()
        );
        println!();
    }

    Ok(aws_config::from_env().region(region_provider).load().await)
}
