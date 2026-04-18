use super_supply_chain;

select
    id,
    name,
    alias,
    unified_social_credit_code,
    target_addr
    from base_companies_infos where name like '%吉安%';

select
    id,
    name,
    alias,
    unified_social_credit_code,
    target_addr
    from base_companies_infos where unified_social_credit_code like '%91130600674681068L%';
