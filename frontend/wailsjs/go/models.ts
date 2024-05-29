export namespace types {
	
	export class DownloadConfig {
	    dir: string;
	
	    static createFrom(source: any = {}) {
	        return new DownloadConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dir = source["dir"];
	    }
	}
	export class S3Config {
	    endpoint: string;
	    accessKey: string;
	    secretKey: string;
	    bucket: string;
	
	    static createFrom(source: any = {}) {
	        return new S3Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.endpoint = source["endpoint"];
	        this.accessKey = source["accessKey"];
	        this.secretKey = source["secretKey"];
	        this.bucket = source["bucket"];
	    }
	}
	export class Config {
	    s3Config: S3Config;
	    downloadConfig: DownloadConfig;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.s3Config = this.convertValues(source["s3Config"], S3Config);
	        this.downloadConfig = this.convertValues(source["downloadConfig"], DownloadConfig);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class File {
	    id: number;
	    path: string;
	    key?: string;
	    isDir: boolean;
	    depth: number;
	    size?: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.path = source["path"];
	        this.key = source["key"];
	        this.isDir = source["isDir"];
	        this.depth = source["depth"];
	        this.size = source["size"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
	    static createFrom(source: any = {}) {
	        return new File(source);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

