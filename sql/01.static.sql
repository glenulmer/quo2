delimiter ###
create or replace procedure plan_deductibles_distinct()
begin
    select distinct is_adult, value from
        (select 1 is_adult, ad_value value from plan_deductibles
         union all
        select 0 is_adult, ch_value from plan_deductibles) z
    order by is_adult, value;
end
###
delimiter ;
