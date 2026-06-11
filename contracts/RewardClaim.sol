pragma solidity ^0.8.0;

import "https://github.com/OpenZeppelin/openzeppelin-contracts/blob/v4.8.0/contracts/security/ReentrancyGuard.sol";
import "https://github.com/OpenZeppelin/openzeppelin-contracts/blob/v4.8.0/contracts/utils/math/SafeMath.sol";

contract RewardClaim is ReentrancyGuard {
    using SafeMath for uint256;

    mapping(address => uint256) public userBalances;
    uint256 public totalRewardPool;

    function claimReward(uint256 amount) public nonReentrant {
        require(amount > 0, "Amount must be greater than zero");
        require(userBalances[msg.sender] >= amount, "Insufficient balance");
        require(totalRewardPool >= amount, "Insufficient reward pool");

        userBalances[msg.sender] = userBalances[msg.sender].sub(amount);
        totalRewardPool = totalRewardPool.sub(amount);

        // Transfer tokens to user
        // ...
    }

    function addUserBalance(address user, uint256 amount) public {
        userBalances[user] = userBalances[user].add(amount);
    }

    function addRewardPool(uint256 amount) public {
        totalRewardPool = totalRewardPool.add(amount);
    }
}