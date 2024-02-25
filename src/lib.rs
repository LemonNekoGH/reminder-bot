use chrono::{DateTime, TimeZone};
use cron_parser::ParseError;

pub mod models;
pub mod schema;

pub fn parse_cron_exp<Tz: TimeZone>(
    exp: &str,
    dt: &DateTime<Tz>,
) -> Result<DateTime<Tz>, ParseError> {
    // check number of expression fields, because cron_parser library won't do this check
    if exp.trim().split(' ').count() < 5 {
        return Err(ParseError::InvalidCron);
    }

    cron_parser::parse(exp, dt)
}
