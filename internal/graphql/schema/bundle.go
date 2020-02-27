package gqlschema

// GraphQL Schema Bundle; auto-created , 2020-02-25 08:22
const schema = `
# Fantom Account type specification
type Account {
    id: ID!
    name: String!
    address: String!
    balance: Amount!
}

# Transaction inside the chain as a result of Transfer
type Transaction {
    id: ID!
    from: Account!
    to: Account!
    amount: Amount!
    timeStamp: Time!
}

# Holds amount of monetary value
scalar Amount

# Holds integere type number for Blockchain indexes
scalar Number

# Holds timestamp
scalar Time

# Raw BlockChain Transaction details
type BlockchainTransaction {
    "Transaction hash identifier."
    hash: ID!

    "Sender address."
    from: String!

    "Recipient address; <null> for contract creation transaction."
    to: String

    "Value transfered in WEI."
    value: Amount!

    "The data send along with the transaction."
    input: String!

    "The number of transactions made by the sender prior to this one."
    nonce: Int!

    "Index of the transaction inside block; <null> for pending."
    txIndex: Int

    "Maximal amount of gas offered by sender for the Transaction."
    gasLimit: Amount!

    "Amount of gas used by Transaction."
    gasUsed: Amount!

    "Price of the gas used by the Transaction."
    gasPrice: Amount!

    "Transaction fee; related to used gas and gas price."
    fee: Amount!

    "Block the Transaction was in; <null> for pending."
    block: BlockchainBlock
}

# Defines pairs of Accounts to be used together
type AccountPair {
    one: Account
    two: Account
}

# Defines input type for Account to Account transfer inside an Account Pair
input TransferInput {
    fromAccountId: ID!
    toAccountId: ID!
    amount: Amount!
}

# Raw BlockChain Block details
type BlockchainBlock {
    "Unique identifier of the Block."
    hash: ID!

    "Number of the Block in the chain."
    number: Number!

    "Timestamp of the Block creation."
    timeStamp: Time!

    "List of hashes of transaction inside the Block."
    txHashes: [String!]!

    "List of transactions inside the Block."
    transactions: [BlockchainTransaction!]!
}

# Root schema definition
schema {
    query: Query
    mutation: Mutation
}

# Entry points for querying the API
type Query {
    "Get an Account either specified by ID, or a random one if ID is not provided."
    account(id:ID):Account!

    "Get list of Accounts specified by their ID."
    accounts(list:[ID!]):[Account!]!

    "Get single pair of accounts, either random ar specified by the ID."
    pair: AccountPair!

    "Get account pairs to be used together."
    pairs: [AccountPair!]!

    "Get raw transaction information for given transaction ID."
    blockchainTransaction(hash:ID!):BlockchainTransaction
}

# data mutation entry points
type Mutation {
    "Transfer funds from one Account to another Account of the same Account Pair."
    transfer(toTransfer: TransferInput!): Transaction

    "Create a burst of transactions from a single Account to random selection of target accounts."
    burst(fromAccountId: ID!, amount: Amount!, targetsCount: Int!): [Transaction!]!
}

`
