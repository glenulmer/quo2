package main

func LoadStaticLookupsOrPanic() {
	App.CustomerLookups = LoadCustomerLookups()
	App.CategOptions = LoadCategOptions()
	App.FilterLookups = LoadFilterLookups()
}
