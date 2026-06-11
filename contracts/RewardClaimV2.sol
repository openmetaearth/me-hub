pragma solidity ^0.8.0;

import "https://github.com/OpenZeppelin/openzeppelin-contracts/blob/v4.8.0/contracts/security/ReentrancyGuard.sol";

contract RewardClaimV2 is ReentrancyGuard {
    mapping(address => uint256) public userBalances;
    uint256 public totalRewardPool;

    function claimReward(uint256 amount) public nonReentrant {
        require(amount > 0, "Amount must be greater than zero");
        require(userBalances[msg.sender] >= amount, "Insufficient balance");
        require(totalRewardPool >= amount, "Insufficient reward pool");

        userBalances[msg.sender] -= amount;
        totalRewardPool -= amount;

        // Transfer tokens to user
        // ...
    }

    function addUserBalance(address user, uint256 amount) public {
        userBalances[user] += amount;
    }

    function addRewardPool(uint256 amount) public {
        totalRewardPool += amount;
    }
}