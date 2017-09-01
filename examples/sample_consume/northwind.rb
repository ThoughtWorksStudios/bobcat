require "json/stream"

BUFFLEN = 16384 # 16k chunks

class ObjectStreamer
  METHODS = %w(start_array end_array start_object end_object key value start_document end_document)
  attr_reader :stack, :keys

  def initialize(parser)
     METHODS.each do |meth_name|
      parser.send(meth_name, &method(meth_name))
    end
  end

  def start_array
    unless top_level?
      stack << []
    end
  end

  def end_array
    unless top_level?
      end_container
    end
  end

  def start_object
    stack << {}
  end

  def end_object
    end_container.tap do |val|
      emit(val) if top_level?
    end
  end

  def end_container
    stack.pop.tap do |val|
      unless top_level?
        case (top = stack[-1])
        when Hash
          top[keys.pop] = val
        when Array
          top << val
        end
      end
    end
  end

  def key(name)
    keys << name
  end

  def value(val)
    case (top = stack[-1])
    when Hash
      top[keys.pop] = val
    when Array
      top << val
    else
      stack << val
    end
  end

  def start_document
    @stack = []
    @keys = []
  end

  def end_document
    stack.pop
    unless stack.empty? && keys.empty?
      raise "parse stack not empty! invalid JSON!"
    end
  end

  def top_level?
    stack.size == 0
  end

  def emit(obj)
    p obj
  end
end

class Sqlizer
  def initialize
    @parser = JSON::Stream::Parser.new
    @streamer = ObjectStreamer.new(@parser)
  end

  def <<(data)
    begin
      @parser << data
    rescue JSON::Stream::ParserError => e
      raise "Failed to parse JSON: #{e.inspect}"
    end
  end
end

sqlizer = Sqlizer.new

if ARGV.size == 0 && STDIN.tty?
  STDERR.puts "You must provide a file to read or pipe input to this script"
  exit 1
end

while buf = ARGF.read(BUFFLEN) do
  sqlizer << buf
end
