// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

contract MedicalAccess {
    struct AccessRequest {
        string recordId;
        string patient;
        string clinic;
        bytes32 tokenHash; // keccak256 hash of token
        uint256 expiresAt;
        string status; // "pending", "approved", "consumed", "revoked"
    }

    mapping(string => AccessRequest) private requests;

    event RequestCreated(string requestId, string recordId, string patient, string clinic);
    event RequestApproved(string requestId, bytes32 tokenHash, uint256 expiresAt);
    event RequestConsumed(string requestId, address consumer);
    event RequestRevoked(string requestId, string reason);

    function requestAccess(
        string calldata requestId,
        string calldata recordId,
        string calldata patient,
        string calldata clinic
    ) external {
        AccessRequest storage ar = requests[requestId];
        // simple guard: if status non-empty - existed
        require(bytes(ar.status).length == 0, "request exists");
        ar.recordId = recordId;
        ar.patient = patient;
        ar.clinic = clinic;
        ar.tokenHash = bytes32(0);
        ar.expiresAt = 0;
        ar.status = "pending";

        emit RequestCreated(requestId, recordId, patient, clinic);
    }

    function approveAccess(
        string calldata requestId,
        bytes32 tokenHash,
        uint256 expiresAt
    ) external {
        AccessRequest storage ar = requests[requestId];
        require(bytes(ar.status).length != 0, "request not found");
        require(keccak256(bytes(ar.status)) == keccak256(bytes("pending")), "not pending");
        ar.tokenHash = tokenHash;
        ar.expiresAt = expiresAt;
        ar.status = "approved";

        emit RequestApproved(requestId, tokenHash, expiresAt);
    }

    function verifyToken(string calldata requestId, bytes32 tokenHash) external view returns (bool) {
        AccessRequest storage ar = requests[requestId];
        if (bytes(ar.status).length == 0) return false;
        if (ar.tokenHash != tokenHash) return false;
        if (block.timestamp > ar.expiresAt) return false;
        if (keccak256(bytes(ar.status)) != keccak256(bytes("approved"))) return false;
        return true;
    }

    // Optional: mark consumed
    function consumeAccess(string calldata requestId) external {
        AccessRequest storage ar = requests[requestId];
        require(bytes(ar.status).length != 0, "request not found");
        require(keccak256(bytes(ar.status)) == keccak256(bytes("approved")), "not approved");
        ar.status = "consumed";
        emit RequestConsumed(requestId, msg.sender);
    }

    function getRequest(string calldata requestId) external view returns (
        string memory recordId,
        string memory patient,
        string memory clinic,
        bytes32 tokenHash,
        uint256 expiresAt,
        string memory status
    ) {
        AccessRequest storage ar = requests[requestId];
        return (ar.recordId, ar.patient, ar.clinic, ar.tokenHash, ar.expiresAt, ar.status);
    }
}