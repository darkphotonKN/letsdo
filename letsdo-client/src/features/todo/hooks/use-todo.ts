import { useMutation, useQueryClient } from "@tanstack/react-query";
import { todoService } from "../services/api";
import type { CreateTodoRequest } from "../types";
import { useToast } from "@/components/ui/use-toast";

const QUERY_KEY = "todos";

export const useCreateTodo = () => {
  const queryClient = useQueryClient();
  const { toast } = useToast();

  return useMutation({
    mutationFn: (data: CreateTodoRequest) => todoService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEY] });
      toast({
        title: "Success",
        description: "Todo created successfully",
      });
    },
    onError: (error: Error) => {
      toast({
        title: "Error",
        description: error.message || "Failed to create todo",
        variant: "destructive",
      });
    },
  });
};
