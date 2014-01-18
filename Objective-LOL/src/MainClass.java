import java.io.File;

import org.objectivelol.lang.LOLError;
import org.objectivelol.lang.LOLSource;
import org.objectivelol.lang.LOLValue;
import org.objectivelol.vm.SourceParser;

public class MainClass {

	/**
	 * @param args
	 * @throws LOLError 
	 */
	public static void main(String[] args) throws LOLError {
		LOLSource e = new SourceParser(new File("test.lol")).parse();
		e.prepareSource();

		LOLValue result = e.getGlobalFunction("FUNC3").execute(null, new LOLValue[0]);
	}

}
