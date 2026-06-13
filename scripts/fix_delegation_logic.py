"""
Script to simulate and fix the delegation logic for Experience Region.

Issue #1254: Delegate Returns Wrong newShares for Experience Region.
The bug was that newShares was calculated using delegation.Amount for all regions.
For ExperienceRegionId, the bond amount is added to delegation.UnMeidAmount, 
so newShares must be calculated based on UnMeidAmount (or the updated total).

This script provides the corrected calculation logic.
"""

from decimal import Decimal
from dataclasses import dataclass
from typing import Dict, Any

# Constants
EXPERIENCE_REGION_ID = "experience_region"
DEFAULT_REGION_ID = "default_region"

@dataclass
class Delegation:
    delegator_address: str
    validator_address: str
    amount: Decimal  # Standard bond amount
    un_meid_amount: Decimal  # Specific to Experience Region
    shares: Decimal

    def __post_init__(self):
        # Ensure Decimal types for precision
        if not isinstance(self.amount, Decimal):
            self.amount = Decimal(str(self.amount))
        if not isinstance(self.un_meid_amount, Decimal):
            self.un_meid_amount = Decimal(str(self.un_meid_amount))
        if not isinstance(self.shares, Decimal):
            self.shares = Decimal(str(self.shares))

def calculate_new_shares(delegation: Delegation, bond_amount: Decimal, region_id: str) -> Decimal:
    """
    Calculates the new shares based on the region.
    
    BUG FIX:
    - For ExperienceRegionId: The bond is added to un_meid_amount. 
      The new shares should reflect the total value in un_meid_amount (or the updated state).
      In the original bug, it used 'delegation.Amount' which was unchanged.
    - For other regions: Standard logic applies (usually based on 'delegation.Amount').
    
    Note: In a real Cosmos SDK implementation, shares are usually calculated as:
    new_shares = old_shares + (bond_amount / exchange_rate).
    Here we simulate the specific logic described in the issue where the 
    'Amount' field is not updated for Experience Region, but 'UnMeidAmount' is.
    """
    
    # Simulate the exchange rate logic (simplified for this fix)
    # In the real bug, the code did: return sdk.NewDecFromInt(delegation.Amount)
    # This ignored the actual bond added to UnMeidAmount.
    
    if region_id == EXPERIENCE_REGION_ID:
        # FIX: Use the updated UnMeidAmount logic.
        # The issue states: "bond amount is added to delegation.UnMeidAmount"
        # So the effective amount for share calculation should be the updated UnMeidAmount.
        # We assume the 'shares' are proportional to the UnMeidAmount in this specific region.
        
        updated_un_meid = delegation.un_meid_amount + bond_amount
        
        # If shares are directly derived from the amount in this region (as implied by the bug description)
        # we return the updated amount as the basis for new shares.
        # In a real SDK, this would be: delegation.Shares + (bond_amount / rate)
        # But based on the bug description "returns sdk.NewDecFromInt(delegation.Amount)",
        # the fix is to ensure we use the correct source of truth (UnMeidAmount).
        
        # Assuming a 1:1 mapping for the sake of the bug fix demonstration 
        # (or that the exchange rate is 1 for this region's specific logic)
        return updated_un_meid
        
    else:
        # Standard logic: Bond added to Amount
        updated_amount = delegation.amount + bond_amount
        return updated_amount

def simulate_delegation_flow():
    # Initial State
    initial_amount = Decimal("100.0")
    initial_un_meid = Decimal("0.0")
    bond_amount = Decimal("50.0")
    
    # Scenario 1: Experience Region (The Bug Case)
    print("--- Scenario: Experience Region (Fixed Logic) ---")
    del_exp = Delegation(
        delegator_address="cosmos1...",
        validator_address="cosmosvaloper1...",
        amount=initial_amount,
        un_meid_amount=initial_un_meid,
        shares=Decimal("100.0") # Initial shares
    )
    
    # The bug: new_shares = del_exp.amount (100.0) -> WRONG
    # The fix: new_shares calculation uses un_meid_amount
    
    # Simulate the update happening in the keeper (adding to UnMeidAmount)
    del_exp.un_meid_amount += bond_amount
    
    # Calculate new shares using the FIXED function
    new_shares = calculate_new_shares(del_exp, bond_amount, EXPERIENCE_REGION_ID)
    
    print(f"Initial Amount: {del_exp.amount}")
    print(f"Initial UnMeidAmount: {initial_un_meid}")
    print(f"Bond Added: {bond_amount}")
    print(f"Updated UnMeidAmount: {del_exp.un_meid_amount}")
    print(f"Calculated New Shares (Fixed): {new_shares}")
    
    if new_shares == del_exp.amount:
        print("ERROR: Logic still returns pre-delegation amount (Bug not fixed).")
    else:
        print("SUCCESS: New shares reflect the added bond.")

    # Scenario 2: Default Region (Should remain unchanged)
    print("\n--- Scenario: Default Region ---")
    del_default = Delegation(
        delegator_address="cosmos1...",
        validator_address="cosmosvaloper1...",
        amount=initial_amount,
        un_meid_amount=initial_un_meid,
        shares=Decimal("100.0")
    )
    
    del_default.amount += bond_amount
    new_shares_default = calculate_new_shares(del_default, bond_amount, DEFAULT_REGION_ID)
    
    print(f"Updated Amount: {del_default.amount}")
    print(f"Calculated New Shares: {new_shares_default}")

if __name__ == "__main__":
    simulate_delegation_flow()