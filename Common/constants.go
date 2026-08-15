package Common

/*
MaxDaysIn3Months is the maximum number of days in 3 months (3*30 days), used for calculating expiration dates and other time-related logic.
*/
const MaxDaysIn3Months = 90

const FullTagPattern = `^[\w][\w.-]{0,127}$`

const TeamNamePattern = `^[a-z][a-z0-9]+$`

const RepositoryNamePattern = `^[a-z0-9]+(?:(?:[._]|__|[-]+)[a-z0-9]+)*(?:/[a-z0-9]+(?:(?:[._]|__|[-]+)[a-z0-9]+)*)*$`

const RobotNamePattern = `^[a-z][a-z0-9_]{1,254}$`

const DefaultTagExpirationSeconds = 1209600 // 14 days in seconds
