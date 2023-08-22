
#[cfg(test)]
mod tests {
    use enclave::kms::parse_raw_data_key;

    #[test]
    fn test_parse_data_key() {
        let input = "CIPHERTEXT: AQIDAHibG56zE+5ETucK/dNFt/HXq5a15TdzUDnATHlDz64jcgHXsP2p8txbiLHl3SXv33kCAAAAfjB8BgkqhkiG9w0BBwagbzBtAgEAMGgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMf2SUP7bhkKZnB9iPAgEQgDtqi3WzX6ozNd9zXjRW6PuPZyEIjvvftssrjRUPQAgxkdvYqcOfHDPU+cjHCmiTKDxh8LSF+6hF5UMf1w==\nPLAINTEXT: lMEpICaR1z7BNcoItPPfXte0BaPs1ONcqe4KpoER+Q4=\n";
        let (ciphertext, plaintext) = parse_raw_data_key(input.to_string());

        assert_eq!(ciphertext, "AQIDAHibG56zE+5ETucK/dNFt/HXq5a15TdzUDnATHlDz64jcgHXsP2p8txbiLHl3SXv33kCAAAAfjB8BgkqhkiG9w0BBwagbzBtAgEAMGgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMf2SUP7bhkKZnB9iPAgEQgDtqi3WzX6ozNd9zXjRW6PuPZyEIjvvftssrjRUPQAgxkdvYqcOfHDPU+cjHCmiTKDxh8LSF+6hF5UMf1w==");
        assert_eq!(plaintext, "lMEpICaR1z7BNcoItPPfXte0BaPs1ONcqe4KpoER+Q4=");
    }
}