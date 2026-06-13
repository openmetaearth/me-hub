#!/usr/bin/env python3
"""
Script to simulate and verify the fix for Issue #1247:
[Bug]: [Medium] [wstaking] RemoveMeidNFT Uses Wrong Store Key Prefix

This script models the storage logic of the MeidNFT contract to demonstrate:
1. The bug: Set writes to 'RegionKey', Remove deletes from 'AccountKey'.
2. The fix: Remove must delete from 'RegionKey' to match Set.

Note: In a real Solidity environment, this logic resides in the .sol file.
This Python script serves as a simulation and verification tool.
"""

import hashlib
from typing import Dict, Set

# Constants simulating the Solidity storage prefixes
# In Solidity, these would be string constants or bytes32 values
MEID_NFT_ACCOUNT_KEY_PREFIX = "meid_nft_account_"
MEID_NFT_REGION_KEY_PREFIX = "meid_nft_region_"

class MeidNFTStorageSimulator:
    def __init__(self):
        # Simulating the state storage (Key -> Value)
        self.storage: Dict[str, str] = {}
        
        # Simulating the regional index (Region -> Set of Accounts)
        # This is the structure that grows indefinitely due to the bug
        self.regional_index: Dict[str, Set[str]] = {}

    def _get_account_key(self, account: str) -> str:
        return f"{MEID_NFT_ACCOUNT_KEY_PREFIX}{account}"

    def _get_region_key(self, region_id: str) -> str:
        return f"{MEID_NFT_REGION_KEY_PREFIX}{region_id}"

    def set_meid_nft(self, account: str, region_id: str, nft_data: str):
        """
        Simulates SetMeidNFT().
        Writes the regional index under MeidNFTRegionKeyPrefix + regionId.
        """
        region_key = self._get_region_key(region_id)
        
        # Initialize the set if it doesn't exist
        if region_key not in self.regional_index:
            self.regional_index[region_key] = set()
        
        # Add the account to the regional index
        self.regional_index[region_key].add(account)
        
        # Store the actual NFT data under the account key (standard practice)
        account_key = self._get_account_key(account)
        self.storage[account_key] = nft_data

    def remove_meid_nft_buggy(self, account: str, region_id: str):
        """
        Simulates the BUGGY RemoveMeidNFT().
        Deletes from MeidNFTAccountKeyPrefix + account.
        This does NOT match the key used in SetMeidNFT for the regional index.
        """
        # BUG: This key is for the account's NFT data, not the regional index entry
        wrong_key = self._get_account_key(account)
        
        # Attempt to delete the regional index entry using the WRONG key
        # This key likely doesn't exist in the regional_index dict, or if it does,
        # it's not the one created by SetMeidNFT.
        if wrong_key in self.regional_index:
            # This block is unreachable for the regional index logic
            self.regional_index[wrong_key].discard(account)
        else:
            # The regional index entry (created with region_key) is never touched
            pass

        # Clean up the account data (this part works, but the index remains)
        if wrong_key in self.storage:
            del self.storage[wrong_key]

    def remove_meid_nft_fixed(self, account: str, region_id: str):
        """
        Simulates the FIXED RemoveMeidNFT().
        Deletes from MeidNFTRegionKeyPrefix + regionId.
        This matches the key used in SetMeidNFT.
        """
        # FIX: Use the correct region key to access the index
        correct_key = self._get_region_key(region_id)
        
        if correct_key in self.regional_index:
            self.regional_index[correct_key].discard(account)
            # Optional: Clean up empty sets to save storage
            if not self.regional_index[correct_key]:
                del self.regional_index[correct_key]

        # Clean up the account data
        account_key = self._get_account_key(account)
        if account_key in self.storage:
            del self.storage[account_key]

    def get_meid_nft_by_account_buggy(self, account: str):
        """
        Simulates GetMeidNFTByAccount with the iterator bug mentioned.
        If the index is not cleaned, phantom entries might appear in queries.
        """
        # In a real scenario, this might iterate over regions to find the account
        # If the index is dirty, it might return false positives or fail to clean up
        found_regions = []
        for region_key, accounts in self.regional_index.items():
            if account in accounts:
                found_regions.append(region_key)
        return found_regions

    def print_state(self, label: str):
        print(f"\n--- State: {label} ---")
        print(f"Storage Keys: {list(self.storage.keys())}")
        print(f"Regional Index Keys: {list(self.regional_index.keys())}")
        print(f"Regional Index Content: {self.regional_index}")
        print(f"Phantom Entries (Account in Index but not in Storage):")
        for region_key, accounts in self.regional_index.items():
            for acc in accounts:
                if self._get_account_key(acc) not in self.storage:
                    print(f"  - Account {acc} in Region {region_key} (Phantom)")

def main():
    print("=== Simulating Issue #1247: Wrong Store Key Prefix ===")
    
    # 1. Setup
    sim = MeidNFTStorageSimulator()
    account = "user_0x123"
    region = "region_A"
    nft_data = "NFT_Data_X"

    # 2. Set MeidNFT (Writes to Region Key)
    print(f"\n1. Setting MeidNFT for {account} in {region}")
    sim.set_meid_nft(account, region, nft_data)
    sim.print_state("After Set")

    # 3. Remove MeidNFT (Buggy Version)
    print(f"\n2. Removing MeidNFT (BUGGY) for {account} in {region}")
    sim.remove_meid_nft_buggy(account, region)
    sim.print_state("After Buggy Remove")
    
    # Check for the bug
    phantom_count = 0
    for region_key, accounts in sim.regional_index.items():
        for acc in accounts:
            if sim._get_account_key(acc) not in sim.storage:
                phantom_count += 1
    
    if phantom_count > 0:
        print(f"\n[BUG DETECTED] Found {phantom_count} phantom entry/entries in regional index!")
        print("The index was not cleaned because the wrong key prefix was used.")
    else:
        print("\n[OK] No phantom entries found.")

    # 4. Reset and Test Fixed Version
    print("\n\n=== Testing the Fix ===")
    sim_fixed = MeidNFTStorageSimulator()
    
    # Set again
    sim_fixed.set_meid_nft(account, region, nft_data)
    sim_fixed.print_state("After Set (Fixed Context)")

    # Remove with Fixed Logic
    print(f"\n3. Removing MeidNFT (FIXED) for {account} in {region}")
    sim_fixed.remove_meid_nft_fixed(account, region)
    sim_fixed.print_state("After Fixed Remove")

    # Verify cleanup
    phantom_count_fixed = 0
    for region_key, accounts in sim_fixed.regional_index.items():
        for acc in accounts:
            if sim_fixed._get_account_key(acc) not in sim_fixed.storage:
                phantom_count_fixed += 1

    if phantom_count_fixed == 0:
        print("\n[SUCCESS] Regional index cleaned correctly. No storage leak.")
    else:
        print(f"\n[ERROR] Still found {phantom_count_fixed} phantom entries.")

if __name__ == "__main__":
    main()