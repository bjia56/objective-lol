package org.objectivelol.vm;

import java.io.File;
import java.util.Collection;
import java.util.HashMap;

import org.objectivelol.lang.LOLError;
import org.objectivelol.lang.LOLFunction;
import org.objectivelol.lang.LOLNative;
import org.objectivelol.lang.LOLSource;
import org.objectivelol.lang.LOLValue;
import org.objectivelol.libs.MATH;
import org.objectivelol.libs.STDIO;

public class RuntimeEnvironment {

	private static RuntimeEnvironment instance = null;
	
	private final HashMap<String, LOLSource> loadedSources = new HashMap<String, LOLSource>();
	private final HashMap<String, LOLNative> nativeFunctions = new HashMap<String, LOLNative>();
	
	private RuntimeEnvironment() throws LOLError {
		if(instance != null) {
			throw new IllegalStateException("Cannot instantiate more than one instance of RuntimeEnvironment");
		}
		
		loadSource("libs" + File.separator + "MATH.lol", "libs" + File.separator + "STDIO.lol");
		loadNative(new MATH(), new STDIO());
	}
	
	public static RuntimeEnvironment getRuntime() throws LOLError {
		if(instance == null) {
			instance  = new RuntimeEnvironment();
		}
		
		return instance;
	}
	
	public void loadSource(String ... fileNames) throws LOLError {
		for(String s : fileNames) {
			File f = new File(s);
			
			SourceParser sp = new SourceParser(f);
			LOLSource result = sp.parse();
			loadedSources.put(result.getName(), result);
		}
	}
	
	public LOLSource getSource(String name) {
		return loadedSources.get(name);
	}
	
	public void loadNative(LOLNative ... natives) {
		for(LOLNative l : natives) {
			nativeFunctions.put(l.getClass().getSimpleName(), l);
		}
	}
	
	public LOLNative getNative(String name) {
		return nativeFunctions.get(name);
	}
	
	public void execute() throws LOLError {
		for(LOLSource s : loadedSources.values()) {
			for(LOLFunction f : s.getGlobalFunctions()) {
				if(f.getName().equals("MAIN")) {
					f.execute(null, (LOLValue[])null);
					return;
				}
			}
		}
	}
	
	public Collection<LOLSource> getLoadedSources() {
		return loadedSources.values();
	}
	
}
