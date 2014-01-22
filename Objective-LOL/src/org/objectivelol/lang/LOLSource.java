package org.objectivelol.lang;

import java.util.Collection;
import java.util.HashMap;

import org.objectivelol.lang.LOLValue.ValueStruct;

public class LOLSource {

	private final String fileName;
	
	private final HashMap<String, ValueStruct> globalVariables;
	private final HashMap<String, LOLFunction> globalFunctions;
	private final HashMap<String, LOLClass> globalClasses;

	public LOLSource(String fileName, HashMap<String, ValueStruct> globalVariables, HashMap<String, LOLFunction> globalFunctions, HashMap<String, LOLClass> globalClasses) {
		this.fileName = fileName;
		this.globalVariables = globalVariables;
		this.globalFunctions = globalFunctions;
		this.globalClasses = globalClasses;
	}

	public ValueStruct getGlobalVariable(String name) {
		return globalVariables.get(name);
	}
	
	public Collection<ValueStruct> getGlobalVariables() {
		return globalVariables.values();
	}

	public LOLFunction getGlobalFunction(String name) {
		return globalFunctions.get(name);
	}
	
	public Collection<LOLFunction> getGlobalFunctions() {
		return globalFunctions.values();
	}
	
	public LOLClass getGlobalClass(String name) throws LOLError {
		return globalClasses.get(name);
	}
	
	public Collection<LOLClass> getGlobalClasses() {
		return globalClasses.values();
	}
	
	public String getName() {
		return fileName;
	}
	
	public void prepareSource() throws LOLError {
		for(LOLClass c : globalClasses.values()) {
			c.prepareClass();
		}
		
		for(LOLFunction func : globalFunctions.values()) {
			func.prepareFunction();
		}
	}
	
}
