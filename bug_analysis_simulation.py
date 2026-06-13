"""
Simulation of the Solvency Check Bug in DoFixedDeposit.

Issue: [Bug] DoFixedDeposit Solvency Check Uses Wrong Formula
Repository: openmetaearth/me-hub
Context: The Go code in x/wstaking/keeper/msg_server_fixed_deposit.go uses a 
dimensionally inconsistent formula for solvency checks.

This Python script simulates the logic to demonstrate the bug:
1. The "Buggy" formula multiplies by regionShare and divides by initAllocationFunds.
2. The "Correct" formula (based on standard interest calculation) should not include these 
   specific scaling factors in the solvency obligation estimate.
3. When regionShare is small, the Buggy formula yields ~0, causing the check to pass 
   even if the treasury is insolvent.
"""

def calculate_interest_obligation_buggy(
    principal: float, 
    rate: float, 
    time: float, 
    region_share: float, 
    init_allocation_funds: float
) -> float:
    """
    Simulates the BUGGY formula found in the Go code.
    Formula: (Principal * Rate * Time) * (regionShare / initAllocationFunds)
    
    This is dimensionally inconsistent because 'regionShare' and 'initAllocationFunds'
    are not part of the standard interest calculation (P * r * t).
    """
    base_interest = principal * rate * time
    
    # The bug: Scaling the obligation by a tiny fraction
    # If region_share is 0.001 and init_allocation is 1000, the factor is 0.000001
    scaling_factor = region_share / init_allocation_funds
    
    estimated_obligation = base_interest * scaling_factor
    return estimated_obligation

def calculate_interest_obligation_correct(
    principal: float, 
    rate: float, 
    time: float
) -> float:
    """
    Simulates the CORRECT formula for interest obligation.
    Formula: Principal * Rate * Time
    
    This represents the actual liability the treasury must cover.
    """
    return principal * rate * time

def check_solvency_buggy(
    treasury_balance: float,
    principal: float,
    rate: float,
    time: float,
    region_share: float,
    init_allocation_funds: float
) -> bool:
    """
    Simulates the solvency check in the Go code using the buggy formula.
    Returns True if the check passes (Treasury >= Estimated Obligation).
    """
    estimated_obligation = calculate_interest_obligation_buggy(
        principal, rate, time, region_share, init_allocation_funds
    )
    
    # The check passes if the treasury has enough to cover the (incorrectly low) estimate
    return treasury_balance >= estimated_obligation

def check_solvency_correct(
    treasury_balance: float,
    principal: float,
    rate: float,
    time: float
) -> bool:
    """
    Simulates the solvency check using the correct formula.
    """
    actual_obligation = calculate_interest_obligation_correct(principal, rate, time)
    return treasury_balance >= actual_obligation

def main():
    print("--- Solvency Check Bug Simulation ---")
    print("Scenario: Small Region (Low region_share)")
    print("-" * 40)

    # Parameters
    treasury_balance = 1000.0  # Treasury has $1000
    principal = 10000.0        # User deposits $10,000
    rate = 0.10                # 10% interest rate
    time = 1.0                 # 1 year
    
    # Bug Trigger: Small region share (e.g., 0.001) and large allocation funds
    region_share = 0.001       
    init_allocation_funds = 1000000.0 

    # Run Buggy Check
    buggy_pass = check_solvency_buggy(
        treasury_balance, principal, rate, time, region_share, init_allocation_funds
    )
    
    # Run Correct Check
    correct_pass = check_solvency_correct(
        treasury_balance, principal, rate, time
    )

    # Calculate values for display
    buggy_obligation = calculate_interest_obligation_buggy(
        principal, rate, time, region_share, init_allocation_funds
    )
    correct_obligation = calculate_interest_obligation_correct(
        principal, rate, time
    )

    print(f"Treasury Balance:      ${treasury_balance:,.2f}")
    print(f"Deposit Principal:     ${principal:,.2f}")
    print(f"Interest Rate:         {rate*100:.1f}%")
    print(f"Region Share:          {region_share}")
    print(f"Init Allocation Funds: {init_allocation_funds:,.2f}")
    print("-" * 40)
    print(f"Correct Obligation:    ${correct_obligation:,.2f} (P * r * t)")
    print(f"Buggy Obligation:      ${buggy_obligation:,.2f} (Scaled by region/allocation)")
    print("-" * 40)
    
    print(f"Correct Solvency Check: {'PASS' if correct_pass else 'FAIL'}")
    print(f"Buggy Solvency Check:   {'PASS' if buggy_pass else 'FAIL'}")
    
    print("-" * 40)
    if buggy_pass and not correct_pass:
        print("❌ BUG DETECTED: The buggy formula allows the deposit to proceed")
        print("   even though the treasury is actually insolvent!")
        print("   The estimated obligation is near zero due to the small region_share.")
    else:
        print("No discrepancy found in this specific scenario.")

if __name__ == "__main__":
    main()