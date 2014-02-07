package org.objectivelol.vm;

import java.io.File;
import java.util.Collection;
import java.util.HashMap;

import org.objectivelol.lang.LOLError;
import org.objectivelol.lang.LOLFunction;
import org.objectivelol.lang.LOLNative;
import org.objectivelol.lang.LOLSource;
import org.objectivelol.libs.MATH;
import org.objectivelol.libs.STDIO;
import org.objectivelol.libs.TIEM;

public class RuntimeEnvironment {

	private static RuntimeEnvironment instance = null;
	
	private final HashMap<String, LOLSource> loadedSources = new HashMap<String, LOLSource>();
	private final HashMap<String, LOLNative> nativeFunctions = new HashMap<String, LOLNative>();
	
	private File execDir = new File(System.getProperty("user.dir"));
	
	private RuntimeEnvironment(File library) throws LOLError {
		if(instance != null) {
			throw new IllegalStateException("Cannot instantiate more than one instance of RuntimeEnvironment");
		}
		
		if(library.isDirectory()) {
			for(File f : library.listFiles()) {
				if(f.isFile()) {
					loadSource(f);
					
					if(f.getName().equals("MATH.lol")) {
						loadNative(new MATH());
					} else if(f.getName().equals("STDIO.lol")) {
						loadNative(new STDIO());
					} else if(f.getName().equals("TIEM.lol")) {
						loadNative(new TIEM());
					}
				}
			}
		}
	}
	
	private RuntimeEnvironment() throws LOLError {
		this(new File("libs"));
	}
	
	public static RuntimeEnvironment getRuntime() throws LOLError {
		if(instance == null) {
			instance  = new RuntimeEnvironment();
		}
		
		return instance;
	}
	
	public static RuntimeEnvironment getRuntime(File library) throws LOLError {
		if(instance == null) {
			instance = new RuntimeEnvironment(library);
		}
		
		return instance;
	}
	
	public void setExecDir(File execDir) {
		this.execDir = execDir.getAbsoluteFile();
	}
	
	public File getExecDir() {
		return execDir;
	}
	
	public void loadSource(File file) throws LOLError {
		SourceParser sp = new SourceParser(file);
		LOLSource result = sp.parse();
		loadedSources.put(result.getName(), result);
	}
	
	public void loadSource(String ... fileNames) throws LOLError {
		for(String s : fileNames) {
			loadSource(new File(s));
		}
	}
	
	public void loadSource(File[] files) throws LOLError {
		for(File f : files) {
			loadSource(f);
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
					f.execute(null);
					return;
				}
			}
		}
	}
	
	public Collection<LOLSource> getLoadedSources() {
		return loadedSources.values();
	}
	
}
